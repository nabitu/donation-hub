package projectstorage

import (
	"context"
	"fmt"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	_type "github.com/isdzulqor/donation-hub/internal/core/type"
	"log"
	"strconv"
	"strings"
	"time"
)

type Storage struct {
	container *model.Container
}

type DatabaseProject struct {
	ID                int64   `db:"id"`
	Name              string  `db:"name"`
	Description       string  `db:"description"`
	ImageURLs         string  `db:"image_urls"`
	TargetAmount      float64 `db:"target_amount"`
	CollectionAmount  float64 `db:"collection_amount"`
	Currency          string  `db:"currency"`
	Status            string  `db:"status"`
	DonorID           int64   `db:"donor_id"`
	DonorUsername     string  `db:"donor_username"`
	RequesterID       int64   `db:"requester_id"`
	RequesterUsername string  `db:"requester_username"`
	RequesterEmail    string  `db:"requester_email"`
	DueAt             int64   `db:"due_at"`
	CreatedAt         int64   `db:"created_at"`
	UpdatedAt         int64   `db:"updated_at"`
}

type DatabaseDonation struct {
	ID            int64  `db:"id"`
	Amount        int64  `db:"amount"`
	Currency      string `db:"currency"`
	Message       string `db:"message"`
	DonorID       int64  `db:"donor_id"`
	DonorUsername string `db:"donor_username"`
	CreatedAt     int64  `db:"created_at"`
}

func New(container *model.Container) *Storage {
	return &Storage{container: container}
}

func (s *Storage) Submit(ctx context.Context, input model.SubmitProjectInput) (projectId int64, err error) {
	query := `INSERT INTO projects (name, description, target_amount, currency, status, requester_id, due_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := s.container.Connection.DB.ExecContext(ctx, query, input.Title, input.Description, input.TargetAmount, input.Currency, _type.PROJECT_NEED_REVIEW, input.UserID, input.DueAt, time.Now().Unix(), time.Now().Unix())

	fmt.Println("ini query submit")
	fmt.Println(query)
	if err != nil {
		fmt.Println("lah error")
		log.Println(err)
		return
	}

	projectId, err = res.LastInsertId()
	if err != nil {
		return
	}

	// Prepare SQL statement for batch insert
	stmt, err := s.container.Connection.DB.PreparexContext(ctx, `
		INSERT INTO project_images (project_id, url)
		VALUES (?, ?)
	`)
	if err != nil {
		return
	}
	defer stmt.Close()

	// Execute the batch insert
	for _, url := range input.ImageURLs {
		_, err = stmt.Exec(projectId, url)
		if err != nil {
			return
		}
	}

	return
}

func (s *Storage) ReviewByAdmin(ctx context.Context, input model.ReviewProjectByAdminInput) (err error) {
	query := `UPDATE projects SET status = ?, updated_at = ? WHERE id = ?`
	_, err = s.container.Connection.DB.ExecContext(ctx, query, input.Status, time.Now().Unix(), input.ProjectId)
	if err != nil {
		err = fmt.Errorf("error updating user: %v", err)
	}

	return
}

func (s *Storage) ListProject(ctx context.Context, input model.ListProjectInput) (o model.ListProjectOutput, err error) {
	var projects []DatabaseProject
	var args []interface{}

	query := `SELECT p.*, IFNULL(GROUP_CONCAT(pi.url), '') AS image_urls, u.id as requester_id, u.username as requester_username, u.email as requester_email
		FROM projects p
		LEFT JOIN project_images pi ON p.id = pi.project_id
		JOIN users u ON u.id = p.requester_id
 		WHERE 1=1
-- 		GROUP BY p.id
		`

	if input.StartTs != 0 {
		query += " AND p.due_at >= ?"
		args = append(args, input.StartTs)
	}

	if input.EndTs != 0 {
		query += " AND p.due_at <= ?"
		args = append(args, input.EndTs)
	}

	if input.LastKey != "" {
		query += " AND p.id > ?"
		args = append(args, input.LastKey)
	}

	if !input.IsAdmin {
		query += " AND NOT p.status = ?"
		args = append(args, _type.PROJECT_NEED_REVIEW)
	}

	query += " GROUP BY p.id"
	if input.Limit != 0 || input.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, input.Limit)
	}

	fmt.Println(args)

	if err = s.container.Connection.DB.SelectContext(ctx, &projects, query, args...); err != nil {
		fmt.Println(err)
		return
	}

	pLists := make([]model.Project, len(projects))
	for i, p := range projects {
		imageUrls := parseImageUrl(projects[i].ImageURLs)
		pLists[i] = model.Project{
			ID:           p.ID,
			Title:        p.Name,
			Description:  p.Description,
			ImageURLs:    imageUrls,
			DueAt:        p.DueAt,
			TargetAmount: p.TargetAmount,
			Currency:     p.Currency,
			Status:       p.Status,
			Requester: model.Requester{
				ID:       p.RequesterID,
				Username: p.RequesterUsername,
				Email:    p.RequesterEmail,
			},
		}
	}

	lastProjectID := ""
	if len(projects) > 0 {
		lastProjectID = fmt.Sprintf("%d", projects[len(projects)-1].ID)
	}
	o.Projects = pLists
	o.LastKey = lastProjectID

	return
}

func (s *Storage) GetProjectById(ctx context.Context, input model.GetProjectByIdInput) (o model.GetProjectByIdOutput, err error) {
	var p DatabaseProject
	query := `
		SELECT
			p.*,
			COALESCE((
				SELECT SUM(d.amount)
				FROM donations d
				WHERE d.project_id = p.id), 0
			) AS collection_amount,
			(
				SELECT GROUP_CONCAT(url)
				FROM project_images pi
				WHERE pi.project_id = p.id
			) AS image_urls,
			u.id AS requester_id,
			u.username AS requester_username,
			u.email AS requester_email
		FROM
			projects p
		JOIN
			users u ON u.id = p.requester_id
		LEFT JOIN
			donations d ON d.project_id = p.id AND p.id = ?
		WHERE
		    p.id = ?
	`
	err = s.container.Connection.DB.GetContext(ctx, &p, query, input.ProjectId, input.ProjectId)

	o.ID = p.ID
	o.Title = p.Name
	o.Description = p.Description
	o.ImageURLs = parseImageUrl(p.ImageURLs)
	fmt.Println(p)
	o.DueAt = p.DueAt
	o.TargetAmount = p.TargetAmount
	o.CollectionAmount = p.CollectionAmount
	o.Currency = p.Currency
	o.Status = p.Status
	o.Requester = model.Requester{
		ID:       p.RequesterID,
		Username: p.RequesterUsername,
		Email:    p.RequesterEmail,
	}

	return
}

func (s *Storage) DonateToProject(ctx context.Context, input model.DonateToProjectInput) (err error) {
	query := `
	INSERT
	INTO donations (project_id, donor_id, message, amount, currency, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	fmt.Println(query)
	_, err = s.container.Connection.DB.ExecContext(ctx, query, input.ProjectId, input.UserID, input.Message, input.Amount, input.Currency, time.Now().Unix())
	if err != nil {
		return
	}

	return
}

func (s *Storage) ListDonationByProjectId(ctx context.Context, input model.ListProjectDonationInput) (output model.ListProjectDonationOutput, err error) {
	var params []interface{}
	var donations []DatabaseDonation

	// write query get data from table donation
	query := `
	SELECT
		d.id as id,
		d.amount as amount,
		d.currency as currency,
		d.message AS message,
		d.donor_id as donor_id,
		u.username AS donor_username,
		d.created_at AS created_at
	FROM
		donations d
	JOIN
		users u ON u.id = d.donor_id
	WHERE
	    d.project_id = ?
	`

	params = append(params, input.ProjectId)
	// Add last_key condition if provided
	if input.LastKey != "" {
		query += " AND d.id > ?"
		params = append(params, input.LastKey)
	}

	if input.Limit > 0 {
		// Add limit for pagination
		query += " ORDER BY d.id ASC LIMIT ?"
		params = append(params, input.Limit)
	}
	err = s.container.Connection.DB.SelectContext(ctx, &donations, query, params...)

	for _, d := range donations {
		output.Donations = append(output.Donations, model.Donation{
			ID:       d.ID,
			Amount:   d.Amount,
			Currency: d.Currency,
			Message:  d.Message,
			Donor: model.Donor{
				ID:       d.DonorID,
				Username: d.DonorUsername,
			},
			CreatedAt: d.CreatedAt,
		})
	}

	if len(donations) > 0 {
		output.LastKey = strconv.FormatInt(donations[len(donations)-1].ID, 10)
	}

	if err != nil {
		return
	}

	return
}

func (s *Storage) HasName(ctx context.Context, name string) (bool, error) {
	query := "select count(*) from projects where name = ?"
	var exists = false
	err := s.container.Connection.DB.GetContext(ctx, &exists, query, name)

	return exists, err
}

func parseImageUrl(urls string) []string {
	imageUrls := strings.Split(urls, ",")

	if len(imageUrls) > 0 {
		if imageUrls[0] == "" {
			imageUrls = nil
		}
	}

	return imageUrls
}
