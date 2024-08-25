package application

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/nguyenthenguyen/docx"
	"golang.org/x/sync/errgroup"
)

const (
	dateReportLayout     = "02/Jan/2006"
	dateTimeReportLayout = "02/Jan/2006 15:04"
	timestampLayout      = "02_01_2006_15_04_05_0000"
)

var tenantsTemplates = map[string]string{
	"sample": "sample_template.docx",
}

type ContentWithAttachment struct {
	Content    string
	Attachment [][]byte
}

type reportService struct {
	reportFolder     string
	ticketService    TicketService
	productService   ProductService
	customerService  CustomerService
	commentService   CommentService
	leadService      LeadService
	tenantService    TenantService
	attachmentBucket domain.AttachmentBucket
}

type ReportService interface {
	GenerateReport(ctx context.Context, crmTicket domain.Ticket) ([]byte, string, error)
}

type ReportData struct {
	CrmTicket domain.Ticket
	Customer  domain.Customer
	Product   domain.Product
	Lead      domain.Lead
	Tenant    domain.Tenant
	Comments  []domain.Comment
}

func NewReportService(
	reportFolder string,
	ticketService TicketService,
	productService ProductService,
	customerService CustomerService,
	commentService CommentService,
	leadService LeadService,
	tenantService TenantService,
	attachmentBucket domain.AttachmentBucket,
) ReportService {
	return &reportService{
		reportFolder:     reportFolder,
		ticketService:    ticketService,
		productService:   productService,
		customerService:  customerService,
		commentService:   commentService,
		leadService:      leadService,
		tenantService:    tenantService,
		attachmentBucket: attachmentBucket,
	}
}

func (s *reportService) GenerateReport(ctx context.Context, crmTicket domain.Ticket) ([]byte, string, error) {
	var memoryDoc bytes.Buffer

	reportData, err := s.getReportData(ctx, crmTicket)
	if err != nil {
		return nil, "", err
	}

	filename, hasTemplate := tenantsTemplates[reportData.Tenant.CompanyName]
	if !hasTemplate {
		return nil, "", fmt.Errorf("no template found for tenant %s", reportData.Tenant.CompanyName)
	}

	err = s.readReportTemplate(ctx, *reportData, filename, &memoryDoc)
	if err != nil {
		return nil, "", err
	}

	return memoryDoc.Bytes(), fmt.Sprintf("%s-%s-%s", reportData.Tenant.CompanyName, crmTicket.ExternalReference, time.Now().Format(timestampLayout)), nil
}

func (s *reportService) getReportData(ctx context.Context, crmTicket domain.Ticket) (*ReportData, error) {
	reportData := &ReportData{
		CrmTicket: crmTicket,
	}

	wg, newCtx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		customer, err := s.customerService.GetByID(newCtx, crmTicket.CustomerID)
		if err != nil {
			return err
		}
		reportData.Customer = *customer
		return nil
	})

	wg.Go(func() error {
		product, err := s.productService.GetProductByID(newCtx, crmTicket.ProductID)
		if err != nil {
			return err
		}
		reportData.Product = *product
		return nil
	})

	wg.Go(func() error {
		comments, err := s.commentService.GetByTicketID(newCtx, crmTicket.TicketID)
		if err != nil {
			return err
		}
		reportData.Comments = comments
		return nil
	})

	wg.Go(func() error {
		lead, err := s.leadService.GetByID(newCtx, crmTicket.LeadID)
		if err != nil {
			return err
		}
		reportData.Lead = *lead
		return nil
	})

	wg.Go(func() error {
		tenant, err := s.tenantService.GetByID(newCtx, crmTicket.TenantID)
		if err != nil {
			return err
		}
		reportData.Tenant = *tenant
		return nil
	})

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return reportData, nil
}

func (s *reportService) readReportTemplate(ctx context.Context, reportData ReportData, filename string, memDoc io.Writer) error {
	filePath := fmt.Sprintf("%s/%s", s.reportFolder, filename)
	file, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return err
	}

	docEdit := file.Editable()
	defer file.Close()

	err = docEdit.Replace("$claim", reportData.CrmTicket.ExternalReference, -1)
	if err != nil {
		return err
	}

	err = docEdit.Replace("$actual_date", time.Now().Format(dateReportLayout), -1)
	if err != nil {
		return err
	}

	err = docEdit.Replace("$client", fmt.Sprintf("%s %s", reportData.Customer.FirstName, reportData.Customer.LastName), -1)
	if err != nil {
		return err
	}

	err = docEdit.Replace("$brand", reportData.Product.Brand, -1)
	if err != nil {
		return err
	}

	err = docEdit.Replace("$summary", reportData.CrmTicket.Subject, -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$lead", fmt.Sprintf("%s %s", reportData.Lead.FirstName, reportData.Lead.LastName), -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$target_date", reportData.CrmTicket.TargetDate.Format(dateReportLayout), -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$document", ParseDocument(reportData.Customer.Document), -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$address", reportData.Customer.ShippingAddress.Address, -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$zip_code", reportData.Customer.ShippingAddress.ZipCode, -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$product", reportData.Product.Name, -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$serial_number", reportData.Product.SerialNumber, -1)
	if err != nil {
		return err
	}

	content := make([]string, 0)
	contentAttachments := make([][]byte, 0)
	comments := make([]string, 0)
	commentsAttachments := make([][]byte, 0)
	resolution := make([]string, 0)
	resolutionAttachments := make([][]byte, 0)
	for _, comment := range reportData.Comments {
		switch comment.CommentType {
		case domain.CONTENT:
			contentAttachments, err = s.downloadFiles(ctx, comment.Attachments)
			if err != nil {
				return err
			}
			content = append(content, fmt.Sprintf("%s - %s", comment.CreatedAt.Format(dateTimeReportLayout), comment.Content))
		case domain.RESOLUTION:
			resolutionAttachments, err = s.downloadFiles(ctx, comment.Attachments)
			if err != nil {
				return err
			}
			resolution = append(resolution, fmt.Sprintf("%s - %s", comment.CreatedAt.Format(dateTimeReportLayout), comment.Content))
		case domain.COMMENT:
			commentsAttachments, err = s.downloadFiles(ctx, comment.Attachments)
			if err != nil {
				return err
			}
			comments = append(comments, fmt.Sprintf("%s - %s", comment.CreatedAt.Format(dateTimeReportLayout), comment.Content))
		}
	}

	err = docEdit.Replace("$content", strings.Join(content, "\r\n"), -1)
	if err != nil {
		return err
	}
	if len(contentAttachments) > 0 {
		err = docEdit.Replace("$image_content", string(contentAttachments[0]), -1)
		if err != nil {
			return err
		}
	}

	err = docEdit.Replace("$comments", strings.Join(comments, "\r\n"), -1)
	if err != nil {
		return err
	}
	if len(commentsAttachments) > 0 {
		err = docEdit.Replace("$image_comment", string(commentsAttachments[0]), -1)
		if err != nil {
			return err
		}
	}

	err = docEdit.Replace("$resolution", strings.Join(resolution, "\r\n"), -1)
	if err != nil {
		return err
	}
	if len(resolutionAttachments) > 0 {
		err = docEdit.Replace("$image_resolution", string(resolutionAttachments[0]), -1)
		if err != nil {
			return err
		}
	}

	return docEdit.Write(memDoc)
}

func (s *reportService) downloadFiles(ctx context.Context, files []domain.Attachment) ([][]byte, error) {
	downloadedFiles := make([][]byte, 0)
	for _, attachment := range files {
		file, err := s.attachmentBucket.Download(ctx, attachment.Key)
		if err != nil {
			return nil, err
		}
		downloadedFiles = append(downloadedFiles, file)
	}

	return downloadedFiles, nil
}
