package usecase

import (
	"errors"
	"github.com/golang/mock/gomock"
	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	pbt "github.com/oOSomnus/transflate/api/generated/translate"
	"github.com/oOSomnus/transflate/internal/task_manager/repository"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"testing"
	"time"
)

func TestTaskUsecaseImpl_ProcessOCRAndTranslate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repository.NewMockUserRepository(ctrl)
	mockOCRClient := service.NewMockOCRClient(ctrl)
	mockTranslateService := service.NewMockTranslateService(ctrl)

	testCases := []struct {
		name        string
		username    string
		fileContent []byte
		lang        string
		mockSetup   func()
		expected    string
		expectError bool
	}{
		{
			name:        "success",
			username:    "testuser",
			fileContent: []byte("filecontent"),
			lang:        "en",
			mockSetup: func() {
				mockOCRClient.EXPECT().ProcessOCR(gomock.Any(), "en").Return(
					&pb.StringListResponse{
						Lines:   []string{"Hello", "World"},
						PageNum: uint32(1),
					}, nil,
				)
				mockUserRepo.EXPECT().DecreaseBalance("testuser", 1).Return(nil)
				mockTranslateService.EXPECT().TranslateText("HelloWorld").Return(
					&pbt.TranslateResult{Lines: "Translated Text"}, nil,
				)
			},
			expected:    "Translated Text",
			expectError: false,
		},
		{
			name:        "ocr error",
			username:    "testuser",
			fileContent: []byte("filecontent"),
			lang:        "en",
			mockSetup: func() {
				mockOCRClient.EXPECT().ProcessOCR(gomock.Any(), "en").Return(nil, errors.New("ocr error"))
			},
			expected:    "",
			expectError: true,
		},
		{
			name:        "decrease balance error",
			username:    "testuser",
			fileContent: []byte("filecontent"),
			lang:        "en",
			mockSetup: func() {
				mockOCRClient.EXPECT().ProcessOCR(gomock.Any(), "en").Return(
					&pb.StringListResponse{
						Lines:   []string{"Hello", "World"},
						PageNum: uint32(1),
					}, nil,
				)
				mockUserRepo.EXPECT().DecreaseBalance("testuser", 1).Return(errors.New("decrease balance error"))
			},
			expected:    "",
			expectError: true,
		},
		{
			name:        "translation error",
			username:    "testuser",
			fileContent: []byte("filecontent"),
			lang:        "en",
			mockSetup: func() {
				mockOCRClient.EXPECT().ProcessOCR(gomock.Any(), "en").Return(
					&pb.StringListResponse{
						Lines:   []string{"Hello", "World"},
						PageNum: uint32(1),
					}, nil,
				)
				mockUserRepo.EXPECT().DecreaseBalance("testuser", 1).Return(nil)
				mockTranslateService.EXPECT().TranslateText("HelloWorld").Return(nil, errors.New("translation error"))
			},
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				tc.mockSetup()
				taskUsecase := &TaskUsecaseImpl{
					ur:   mockUserRepo,
					ocrc: mockOCRClient,
					ts:   mockTranslateService,
				}
				result, err := taskUsecase.ProcessOCRAndTranslate(tc.username, tc.fileContent, tc.lang)
				if tc.expectError && err == nil {
					t.Errorf("expected error but got none")
				}
				if !tc.expectError && err != nil {
					t.Errorf("did not expect error but got: %v", err)
				}
				if result != tc.expected {
					t.Errorf("expected: %s, got: %s", tc.expected, result)
				}
			},
		)
	}
}

func TestTaskUsecaseImpl_CreateDownloadLinkWithMdString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3Storage := service.NewMockS3StorageService(ctrl)

	testCases := []struct {
		name        string
		mdString    string
		mockSetup   func()
		expected    string
		expectError bool
	}{
		{
			name:     "success",
			mdString: "markdown content",
			mockSetup: func() {
				mockS3Storage.EXPECT().UploadFileToS3(gomock.Any(), gomock.Any(), gomock.Any(), 1).Return(nil)
				mockS3Storage.EXPECT().GeneratePresignedURL(
					gomock.Any(), gomock.Any(), time.Duration(presignedURLExpiry),
				).Return("http://downloadlink.com", nil)
			},
			expected:    "http://downloadlink.com",
			expectError: false,
		},
		{
			name:     "upload error",
			mdString: "markdown content",
			mockSetup: func() {
				mockS3Storage.EXPECT().UploadFileToS3(
					gomock.Any(), gomock.Any(), gomock.Any(), 1,
				).Return(errors.New("upload error"))
			},
			expected:    "",
			expectError: true,
		},
		{
			name:     "generate link error",
			mdString: "markdown content",
			mockSetup: func() {
				mockS3Storage.EXPECT().UploadFileToS3(gomock.Any(), gomock.Any(), gomock.Any(), 1).Return(nil)
				mockS3Storage.EXPECT().GeneratePresignedURL(
					gomock.Any(), gomock.Any(), time.Duration(presignedURLExpiry),
				).Return("", errors.New("generate link error"))
			},
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				tc.mockSetup()
				taskUsecase := &TaskUsecaseImpl{
					s3s: mockS3Storage,
				}
				result, err := taskUsecase.CreateDownloadLinkWithMdString(tc.mdString)
				if tc.expectError && err == nil {
					t.Errorf("expected error but got none")
				}
				if !tc.expectError && err != nil {
					t.Errorf("did not expect error but got: %v", err)
				}
				if result != tc.expected {
					t.Errorf("expected: %s, got: %s", tc.expected, result)
				}
			},
		)
	}
}
