package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin" 
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/repository"
	"github.com/Prototype-1/xtrace/pkg/utils"
	 "github.com/jung-kurt/gofpdf"
)

type InvoiceHandler struct {
	UserRepository    repository.UserRepository
	InvoiceRepository repository.InvoiceRepository
}

func NewInvoiceHandler(userRepo repository.UserRepository, invoiceRepo repository.InvoiceRepository) *InvoiceHandler {
	return &InvoiceHandler{
		UserRepository:    userRepo,
		InvoiceRepository: invoiceRepo,
	}
}

func (h *InvoiceHandler) GenerateInvoicePDF(invoice *models.Invoice) (string, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)
    
    pdf.Cell(40, 10, "Invoice Details")
    pdf.Ln(12) 

    pdf.SetFont("Arial", "", 12)
    pdf.Cell(40, 10, fmt.Sprintf("Invoice ID: %d", invoice.InvoiceID))
    pdf.Ln(8)
    pdf.Cell(40, 10, fmt.Sprintf("Original Amount: %.2f", invoice.OriginalAmount))
    pdf.Ln(8)
    pdf.Cell(40, 10, fmt.Sprintf("Discount: %.2f", invoice.DiscountAmount))
    pdf.Ln(8)
    pdf.Cell(40, 10, fmt.Sprintf("Amount Due: %.2f", invoice.Amount))
    pdf.Ln(8)
    pdf.Cell(40, 10, fmt.Sprint("Payment Type: ", invoice.PaymentType))
    pdf.Ln(8)
    pdf.Cell(40, 10, fmt.Sprintf("Status: %s", invoice.Status))

    // Save the PDF to a temporary file
    fileName := fmt.Sprintf("invoice_%d.pdf", invoice.InvoiceID)
    err := pdf.OutputFileAndClose(fileName)
    if err != nil {
        return "", err
    }

    return fileName, nil
}

func (h *InvoiceHandler) GetUserEmail(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("userID")) 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UserRepository.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user email"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": user.Email}) 
}

func (h *InvoiceHandler) SendInvoice(c *gin.Context) {
    var req struct {
        UserID uint   `json:"userID"`
        Email  string `json:"email"`
    }

    if err := c.ShouldBindJSON(&req); err != nil { 
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    invoice, err := h.InvoiceRepository.GetInvoiceByUserID(req.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoice"})
        return
    }
    if invoice == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
        return
    }

    pdfFile, err := h.GenerateInvoicePDF(invoice)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
        return
    }

    if err = utils.SendEmailWithAttachment(req.Email, "Your Invoice", "Please find the attached invoice.", pdfFile); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send invoice email"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Invoice sent successfully!"})
}


