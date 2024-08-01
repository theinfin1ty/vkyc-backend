package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"os"
	"strconv"
	"strings"

	"time"
	"vkyc-backend/configs"
	"vkyc-backend/models"
	"vkyc-backend/types"
	"vkyc-backend/utils"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/gin-gonic/gin"
)

func CompleteKycProcess(c *gin.Context) {
	var meeting models.Meeting
	meetingCode := c.Param("meetingCode")

	result := configs.DB.First(&meeting, models.Meeting{
		MeetingCode: meetingCode,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting not found") {
		return
	}

	status := utils.ApprovedStatus

	meeting.KycStatus = &status

	result.Save(meeting)

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "KYC process completed successfully",
	}))
}

func GetKYCQuestions(c *gin.Context) {
	var questions []models.Question
	// var questionsListConfig models.Configuration
	// var questionsCountConfig models.Configuration

	// result := configs.DB.Find(&questionsCountConfig, models.Configuration{
	// 	Key: utils.QuestionsCount,
	// })

	// if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Questions Count Configuration not found") {
	// 	return
	// }

	// result = configs.DB.Find(&questionsListConfig, models.Configuration{
	// 	Key: utils.QuestionsList,
	// })

	// if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Questions List Configuration not found") {
	// 	return
	// }

	// questionsCount, err := strconv.Atoi(questionsCountConfig.Value)

	// if utils.CheckError(err, http.StatusInternalServerError, c, "Questions Count Configuration Value is Invalid") {
	// 	return
	// }

	// questionsList := strings.Split(questionsListConfig.Value, ",")
	// questionIds := make([]int, questionsCount)

	// for i := 0; i < questionsCount; i++ {
	// 	questionsList[i] = strings.TrimSpace(questionsList[i])
	// 	id, err := strconv.Atoi(questionsList[i])

	// 	if err != nil {
	// 		continue
	// 	}

	// 	questionIds = append(questionIds, id)
	// }

	// result = configs.DB.Where("id IN ?", questionIds).Find(&questions)

	// if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Questions not found") {
	// 	return
	// }

	result := configs.DB.Find(&questions, models.Question{
		IsDeleted: false,
		Status:    utils.ActiveStatus,
	})

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Questions not found") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"questions": questions,
	}))
}

func SubmitAnswer(c *gin.Context) {
	var meeting models.Meeting
	var question models.Question
	var answer models.Answer
	var body types.AnswerInput

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	result := configs.DB.First(&meeting, models.Meeting{
		MeetingCode: body.MeetingCode,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting not found") {
		return
	}

	result = configs.DB.First(&question, models.Question{
		ID: body.QuestionID,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Question not found") {
		return
	}

	configs.DB.First(&answer, models.Answer{
		MeetingID:  meeting.ID,
		QuestionID: question.ID,
	})

	answer.MeetingID = meeting.ID
	answer.QuestionID = question.ID
	answer.Answer = body.Answer

	result = configs.DB.Save(&answer)

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"answer": answer,
	}))
}

func DocumentVerification(c *gin.Context) {
	var meeting models.Meeting
	var document interface{}
	// var targetImageBase64, sourceImageBase64 string
	var err error

	countryCode := c.Request.FormValue("countryCode")
	cardCode := c.Request.FormValue("cardCode")
	verificationType := c.Request.FormValue("verificationType")
	meetingCode := c.Request.FormValue("meetingCode")
	documentName := c.Request.FormValue("documentName")
	verificationId := c.Request.FormValue("verificationId")
	sourceImageBase64 := c.Request.FormValue("sourceImage")
	targetImageBase64 := c.Request.FormValue("targetImage")

	// sourceImage, err := c.FormFile("sourceImage")
	// if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
	// 	return
	// }

	// sourceImageFile, err := sourceImage.Open()
	// if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
	// 	return
	// }

	// defer sourceImageFile.Close()

	// sourceImageBytes, err := io.ReadAll(sourceImageFile)
	// if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
	// 	return
	// }

	// sourceImageBase64 = base64.StdEncoding.EncodeToString(sourceImageBytes)

	// if verificationType == "facematch" {
	// 	targetImage, err := c.FormFile("targetImage")
	// 	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
	// 		return
	// 	}

	// 	targetImageFile, err := targetImage.Open()
	// 	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
	// 		return
	// 	}

	// 	defer targetImageFile.Close()

	// 	targetImageBytes, err := io.ReadAll(targetImageFile)
	// 	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
	// 		return
	// 	}

	// 	targetImageBase64 = base64.StdEncoding.EncodeToString(targetImageBytes)
	// }

	result := configs.DB.First(&meeting, models.Meeting{
		MeetingCode: meetingCode,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting not found") {
		return
	}

	switch verificationType {
	case "ocr":
		document, err = verifyOCR(countryCode, cardCode, sourceImageBase64, &meeting)
	case "liveness":
		document, err = verifyLiveness(sourceImageBase64, &meeting, verificationId)
	case "facematch":
		document, err = verifyFaceMatch(sourceImageBase64, targetImageBase64, &meeting, verificationId)
	case "signature":
		document, err = saveDocument(sourceImageBase64, &meeting, verificationType, nil)
	case "document":
		document, err = saveDocument(sourceImageBase64, &meeting, verificationType, &documentName)
	}

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"document": document,
	}))
}

func saveDocument(sourceImage string, meeting *models.Meeting, verificationType string, documentName *string) (*models.Document, error) {
	imagePath := fmt.Sprintf("uploads/%s/documents", meeting.MeetingCode)
	imageName := fmt.Sprintf("%s-%d.jpg", verificationType, time.Now().UnixMilli())
	fullPath := fmt.Sprintf("%s/%s", imagePath, imageName)

	document := models.Document{
		Type:      verificationType,
		Image:     fullPath,
		MeetingID: meeting.ID,
		// APIResponse: nil,
		APIStatus: nil,
		Score:     nil,
		Name:      documentName,
	}

	err := addTimestampToImageAndSave(sourceImage, imagePath, fullPath, *meeting)

	if err != nil {
		return nil, err
	}

	result := configs.DB.Create(&document)

	if result.Error != nil {
		return nil, result.Error
	}

	err = deleteExtraDocuments(meeting.ID, 10, verificationType)

	if err != nil {
		return nil, err
	}

	return &document, nil
}

func verifyOCR(countryCode string, cardCode string, sourceImage string, meeting *models.Meeting) (*types.OCRVerificationResponse, error) {
	var ocrResponse types.OCRResponse
	var rawOCRResponse types.OCRResponse
	var apiConfiguration, apiKeyConfiguration models.Configuration

	result := configs.DB.First(&apiConfiguration, models.Configuration{
		Key:    utils.OCRAPI,
		Status: utils.ActiveStatus,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	if apiConfiguration.Value == "" {
		return nil, errors.New("ocr api not found in configurations")
	}

	result = configs.DB.First(&apiKeyConfiguration, models.Configuration{
		Key:    utils.OCRAPIKey,
		Status: utils.ActiveStatus,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	if apiKeyConfiguration.Value == "" {
		return nil, errors.New("ocr api key not found in configurations")
	}

	apiUrl := apiConfiguration.Value
	apiKey := apiKeyConfiguration.Value

	data := url.Values{}
	data.Add("country_code", countryCode)
	data.Add("card_code", cardCode)
	data.Add("scan_image_base64", sourceImage)

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Api-Key", apiKey)

	res, err := utils.Client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &ocrResponse)

	if err != nil {
		return nil, errors.New("unable to process ocr response")
	}

	err = json.Unmarshal(body, &rawOCRResponse)

	if err != nil {
		return nil, errors.New("unable to process ocr response")
	}

	if ocrResponse.Data != nil && ocrResponse.Data.OCRData != nil && len(*ocrResponse.Data.OCRData) != 0 {
		delete(*ocrResponse.Data.OCRData, "face")
	}

	if ocrResponse.Data != nil && ocrResponse.Data.MRZData != nil && len(*ocrResponse.Data.MRZData) != 0 {
		delete(*ocrResponse.Data.MRZData, "face_image")
	}

	ocrResponseJson, err := json.Marshal(ocrResponse)

	if err != nil {
		return nil, err
	}

	encryptedOCRResponseJson, err := utils.EncryptData(ocrResponseJson)

	if err != nil {
		return nil, err
	}

	encryptedOCRResponseJsonString := base64.StdEncoding.EncodeToString(encryptedOCRResponseJson)

	imagePath := fmt.Sprintf("uploads/%s/documents", meeting.MeetingCode)
	imageName := fmt.Sprintf("ocr-%d.jpg", time.Now().UnixMilli())
	fullPath := fmt.Sprintf("%s/%s", imagePath, imageName)

	document := models.Document{
		Type:        "ocr",
		APIResponse: &encryptedOCRResponseJsonString,
		APIStatus:   &ocrResponse.Status,
		Image:       fullPath,
		MeetingID:   meeting.ID,
	}

	err = addTimestampToImageAndSave(sourceImage, imagePath, fullPath, *meeting)

	if err != nil {
		return nil, err
	}

	result = configs.DB.Create(&document)

	if result.Error != nil {
		return nil, result.Error
	}

	err = deleteExtraDocuments(meeting.ID, 10, "ocr")

	if err != nil {
		return nil, err
	}

	return &types.OCRVerificationResponse{
		Document:    document,
		OCRResponse: rawOCRResponse,
	}, nil
}

func verifyLiveness(sourceImage string, meeting *models.Meeting, verificationId string) (*models.Document, error) {
	var livenessResponse types.LivenessResponse
	var apiConfiguration,
		apiKeyConfiguration,
		passScoreConfiguration models.Configuration

	result := configs.DB.First(&apiConfiguration, models.Configuration{
		Key:    utils.LivenessAPI,
		Status: utils.ActiveStatus,
	})

	if result.Error != nil {
		return nil, errors.New("liveness api not found in configurations")
	}

	result = configs.DB.First(&apiKeyConfiguration, models.Configuration{
		Key:    utils.LivenessAPIKey,
		Status: utils.ActiveStatus,
	})

	if result.Error != nil {
		return nil, errors.New("liveness api key not found in configurations")
	}

	result = configs.DB.First(&passScoreConfiguration, models.Configuration{
		Key:    utils.MinLivenessScore,
		Status: utils.ActiveStatus,
	})

	if result.Error != nil {
		return nil, errors.New("liveness pass score not found in configurations")
	}

	if apiConfiguration.Value == "" {
		return nil, errors.New("liveness api not found in configurations")
	}

	if apiKeyConfiguration.Value == "" {
		return nil, errors.New("liveness api key not found in configurations")
	}

	if passScoreConfiguration.Value == "" {
		return nil, errors.New("liveness pass score not found in configurations")
	}

	passScore, err := strconv.ParseFloat(passScoreConfiguration.Value, 64)

	if err != nil {
		return nil, err
	}

	apiUrl := apiConfiguration.Value
	apiKey := apiKeyConfiguration.Value

	imageData, err := base64.StdEncoding.DecodeString(sourceImage)

	if err != nil {
		return nil, err
	}

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	part, err := writer.CreateFormFile("liveness_image", "liveness.jpg")

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(imageData))

	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiUrl, data)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Api-Key", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := utils.Client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &livenessResponse)

	if err != nil {
		return nil, err
	}

	var score float64

	if livenessResponse.Data != nil {
		score = (*livenessResponse.Data)["score"].(float64)
	}

	apiResponse := string(body)

	imagePath := fmt.Sprintf("uploads/%s/documents", meeting.MeetingCode)
	imageName := fmt.Sprintf("liveness-%d.jpg", time.Now().UnixMilli())
	fullPath := fmt.Sprintf("%s/%s", imagePath, imageName)
	apiStatus := utils.FailStatus

	if score >= passScore {
		apiStatus = utils.PassStatus
	}

	document := models.Document{
		Type:           "liveness",
		APIResponse:    &apiResponse,
		APIStatus:      &apiStatus,
		Image:          fullPath,
		MeetingID:      meeting.ID,
		Score:          &score,
		VerificationID: &verificationId,
	}

	err = addTimestampToImageAndSave(sourceImage, imagePath, fullPath, *meeting)

	if err != nil {
		return nil, err
	}

	result = configs.DB.Create(&document)

	if result.Error != nil {
		return nil, result.Error
	}

	err = deleteExtraDocuments(meeting.ID, 10, "liveness")

	if err != nil {
		return nil, err
	}

	return &document, nil
}
func verifyFaceMatch(sourceImage string, targetImage string, meeting *models.Meeting, verificationId string) (*models.Document, error) {
	var facematchResponse types.FacematchResponse
	var apiConfiguration,
		apiKeyConfiguration,
		passScoreConfiguration models.Configuration

	result := configs.DB.First(&apiConfiguration, models.Configuration{
		Key:    utils.FacematchAPI,
		Status: utils.ActiveStatus,
	})

	if result.Error != nil {
		return nil, errors.New("facematch api not found in configurations")
	}

	result = configs.DB.First(&apiKeyConfiguration, models.Configuration{
		Key:    utils.FacematchAPIKey,
		Status: utils.ActiveStatus,
	})

	if result.Error != nil {
		return nil, errors.New("facematch api key not found in configurations")
	}

	result = configs.DB.First(&passScoreConfiguration, models.Configuration{
		Key:    utils.MinFacematchScore,
		Status: utils.ActiveStatus,
	})

	if result.Error != nil {
		return nil, errors.New("facematch pass score not found in configurations")
	}

	if apiConfiguration.Value == "" {
		return nil, errors.New("facematch api not found in configurations")
	}

	if apiKeyConfiguration.Value == "" {
		return nil, errors.New("facematch api key not found in configurations")
	}

	if passScoreConfiguration.Value == "" {
		return nil, errors.New("facematch pass score not found in configurations")
	}

	passScore, err := strconv.ParseFloat(passScoreConfiguration.Value, 64)

	if err != nil {
		return nil, err
	}

	apiUrl := apiConfiguration.Value
	apiKey := apiKeyConfiguration.Value

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	// Add source image
	sourceImageField, err := writer.CreateFormField("source_image_base64")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(sourceImageField, bytes.NewBufferString(sourceImage))
	if err != nil {
		return nil, err
	}

	// Add target image
	targetImageField, err := writer.CreateFormField("target_image_base64")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(targetImageField, bytes.NewBufferString(targetImage))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiUrl, data)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Api-Key", apiKey)

	res, err := utils.Client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &facematchResponse)

	if err != nil {
		return nil, err
	}

	var score float64

	var facematchResponseJsonString string

	if facematchResponse.Data != nil {
		score = (facematchResponse.Data)["score"].(float64)
		delete(facematchResponse.Data, "retimg1")
		delete(facematchResponse.Data, "retimg2")

		facematchResponseJson, err := json.Marshal(facematchResponse)

		if err != nil {
			return nil, err
		}

		facematchResponseJsonString = string(facematchResponseJson)
	}

	apiStatus := utils.FailStatus
	if score*100 >= passScore {
		apiStatus = utils.PassStatus
	}

	imagePath := fmt.Sprintf("uploads/%s/documents", meeting.MeetingCode)
	imageName := fmt.Sprintf("facematch-%d.jpg", time.Now().UnixMilli())
	fullPath := fmt.Sprintf("%s/%s", imagePath, imageName)

	document := models.Document{
		Type:           "facematch",
		APIResponse:    &facematchResponseJsonString,
		APIStatus:      &apiStatus,
		Image:          fullPath,
		MeetingID:      meeting.ID,
		Score:          &score,
		VerificationID: &verificationId,
	}

	err = addTimestampToImageAndSave(sourceImage, imagePath, fullPath, *meeting)

	if err != nil {
		return nil, err
	}

	result = configs.DB.Create(&document)

	if result.Error != nil {
		return nil, result.Error
	}

	err = deleteExtraDocuments(meeting.ID, 10, "facematch")

	if err != nil {
		return nil, err
	}

	return &document, nil
}

func deleteExtraDocuments(meetingId uint, limit int, documentType string) error {
	var documents []models.Document

	result := configs.DB.Order("id desc").Find(&documents, models.Document{
		Type:      documentType,
		MeetingID: meetingId,
	})

	if result.Error != nil {
		return result.Error
	}

	if len(documents) > limit {
		extraDocuments := documents[limit:]

		for _, document := range extraDocuments {
			err := os.Remove(document.Image)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			result = configs.DB.Delete(&document)
			if result.Error != nil {
				fmt.Println(result.Error.Error())
				continue
			}
		}
	}

	return nil
}

func addTimestampToImageAndSave(base64Image string, filePath string, fullPath string, meeting models.Meeting) error {
	timestamp := time.Now().Format(utils.DateLayout)

	if meeting.Location == nil {
		locStr := ""
		meeting.Location = &locStr
	}

	watermark := fmt.Sprintf("%s %s", timestamp, *meeting.Location)

	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		fmt.Println("Error decoding base64 string:", err)
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		fmt.Println("Error creating image from decoded data:", err)
		return err
	}

	// Create a new RGBA image
	rgba := image.NewRGBA(img.Bounds())

	// Draw the original image on the new RGBA image
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Add text at the bottom-right corner
	textX := rgba.Bounds().Dx() - len(watermark)*9
	textY := rgba.Bounds().Dy() - 30
	col := image.White
	point := fixed.P(textX, textY)
	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(watermark)

	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, rgba, nil)

	if err != nil {
		fmt.Println("Error encoding image:", err)
		return err
	}

	timestampedImageData := buffer.Bytes()

	encryptedTimestampedImageData, err := utils.EncryptData(timestampedImageData)

	if err != nil {
		fmt.Println("Error encrypting timestamped image data:", err)
		return err
	}

	err = os.MkdirAll(filePath, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	err = os.WriteFile(fullPath, encryptedTimestampedImageData, 0644)

	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	return nil
}
