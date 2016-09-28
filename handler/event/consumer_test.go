package event

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"testing"
)

const (
	queueUrl       = "url://test"
	receiptHandle  = "receiptHandle"
	receiptHandle2 = "receiptHandle2"
	messageBody    = "messageBody"
	messageBody2   = "messageBody2"
)

type consumerMockContext struct {
	mockCtrl             *gomock.Controller
	sqsClient            *mocks.MockSQSAPI
	processor            *mocks.MockProcessor
	getQueueUrlInput     *sqs.GetQueueUrlInput
	getQueueUrlOutput    *sqs.GetQueueUrlOutput
	receiveMessageInput  *sqs.ReceiveMessageInput
	receiveMessageOutput *sqs.ReceiveMessageOutput
	sqsMessage           *sqs.Message
	sqsMessage2          *sqs.Message
	deleteMessageInput   *sqs.DeleteMessageInput
	deleteMessageInput2  *sqs.DeleteMessageInput
}

func NewConsumerMockContext(t *testing.T) *consumerMockContext {
	context := consumerMockContext{}
	context.mockCtrl = gomock.NewController(t)
	context.sqsClient = mocks.NewMockSQSAPI(context.mockCtrl)
	context.processor = mocks.NewMockProcessor(context.mockCtrl)

	context.sqsMessage = &sqs.Message{
		Body:          aws.String(messageBody),
		ReceiptHandle: aws.String(receiptHandle),
	}

	context.sqsMessage2 = &sqs.Message{
		Body:          aws.String(messageBody2),
		ReceiptHandle: aws.String(receiptHandle2),
	}

	context.getQueueUrlInput = &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	context.getQueueUrlOutput = &sqs.GetQueueUrlOutput{
		QueueUrl: aws.String(queueUrl),
	}

	context.receiveMessageInput = &sqs.ReceiveMessageInput{
		QueueUrl:          aws.String(queueUrl),
		VisibilityTimeout: aws.Int64(visibilityTimeout),
		WaitTimeSeconds:   aws.Int64(waitTimeSeconds),
	}

	context.receiveMessageOutput = &sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{context.sqsMessage, context.sqsMessage2},
	}

	context.deleteMessageInput = &sqs.DeleteMessageInput{
		ReceiptHandle: aws.String(receiptHandle),
		QueueUrl:      aws.String(queueUrl),
	}

	context.deleteMessageInput2 = &sqs.DeleteMessageInput{
		ReceiptHandle: aws.String(receiptHandle2),
		QueueUrl:      aws.String(queueUrl),
	}

	return &context
}

func TestNewConsumerNilSQS(t *testing.T) {
	context := NewConsumerMockContext(t)
	defer context.mockCtrl.Finish()

	_, err := NewConsumer(nil, context.processor)
	if err == nil {
		t.Error("Expected an error when sqs is nil")
	}
}

func TestNewConsumerNilProcessor(t *testing.T) {
	context := NewConsumerMockContext(t)
	defer context.mockCtrl.Finish()

	_, err := NewConsumer(context.sqsClient, nil)
	if err == nil {
		t.Error("Expected an error when processor is nil")
	}
}

func TestNewConsumerGetQueueUrlFails(t *testing.T) {
	context := NewConsumerMockContext(t)
	defer context.mockCtrl.Finish()

	context.sqsClient.EXPECT().GetQueueUrl(gomock.Eq(context.getQueueUrlInput)).Return(nil, errors.New(""))

	_, err := NewConsumer(context.sqsClient, context.processor)

	if err == nil {
		t.Error("Expected an error when getQueueUrl fails")
	}
}

func TestNewConsumerGetQueueUrlIsEmpty(t *testing.T) {
	context := NewConsumerMockContext(t)
	defer context.mockCtrl.Finish()

	context.sqsClient.EXPECT().GetQueueUrl(gomock.Eq(context.getQueueUrlInput)).Return(&sqs.GetQueueUrlOutput{}, nil)

	_, err := NewConsumer(context.sqsClient, context.processor)

	if err == nil {
		t.Error("Expected an error when getQueueUrl output is empty")
	}
}

func TestNewConsumer(t *testing.T) {
	context := NewConsumerMockContext(t)
	defer context.mockCtrl.Finish()

	context.sqsClient.EXPECT().GetQueueUrl(gomock.Eq(context.getQueueUrlInput)).Return(context.getQueueUrlOutput, nil)

	c, err := NewConsumer(context.sqsClient, context.processor)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	if c == nil {
		t.Error("Consumer should not be nil")
	}
}

func TestPollForEventsReceiveMessageFails(t *testing.T) {
	mockContext := NewConsumerMockContext(t)
	defer mockContext.mockCtrl.Finish()

	mockContext.sqsClient.EXPECT().GetQueueUrl(gomock.Eq(mockContext.getQueueUrlInput)).Return(mockContext.getQueueUrlOutput, nil)

	c, err := NewConsumer(mockContext.sqsClient, mockContext.processor)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	pollCount := 0

	mockContext.sqsClient.EXPECT().ReceiveMessage(mockContext.receiveMessageInput).
		Return(nil, errors.New("Receive message fails")).Do(func(x interface{}) {
		pollCount++
		if pollCount == 1 {
			cancel()
		}
	})

	c.PollForEvents(ctx)
}

func TestPollForEventsReceiveMessageOutputNil(t *testing.T) {
	mockContext := NewConsumerMockContext(t)
	defer mockContext.mockCtrl.Finish()

	mockContext.sqsClient.EXPECT().GetQueueUrl(gomock.Eq(mockContext.getQueueUrlInput)).Return(mockContext.getQueueUrlOutput, nil)

	c, err := NewConsumer(mockContext.sqsClient, mockContext.processor)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	pollCount := 0

	mockContext.sqsClient.EXPECT().ReceiveMessage(mockContext.receiveMessageInput).Return(nil, nil).Do(func(x interface{}) {
		pollCount++
		if pollCount == 1 {
			cancel()
		}
	})

	c.PollForEvents(ctx)
}

func TestPollForEventsReceiveMessageOutputMessagesNil(t *testing.T) {
	mockContext := NewConsumerMockContext(t)
	defer mockContext.mockCtrl.Finish()

	mockContext.sqsClient.EXPECT().GetQueueUrl(gomock.Eq(mockContext.getQueueUrlInput)).Return(mockContext.getQueueUrlOutput, nil)

	c, err := NewConsumer(mockContext.sqsClient, mockContext.processor)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	pollCount := 0

	receiveMessageOutput := &sqs.ReceiveMessageOutput{}
	mockContext.sqsClient.EXPECT().ReceiveMessage(mockContext.receiveMessageInput).
		Return(receiveMessageOutput, nil).Do(func(x interface{}) {
		pollCount++
		if pollCount == 1 {
			cancel()
		}
	})

	c.PollForEvents(ctx)
}

func TestPollForEventsFirstProcessEventFails(t *testing.T) {
	mockContext := NewConsumerMockContext(t)
	defer mockContext.mockCtrl.Finish()

	mockContext.sqsClient.EXPECT().GetQueueUrl(gomock.Eq(mockContext.getQueueUrlInput)).Return(mockContext.getQueueUrlOutput, nil)

	c, err := NewConsumer(mockContext.sqsClient, mockContext.processor)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	pollCount := 0

	mockContext.sqsClient.EXPECT().ReceiveMessage(mockContext.receiveMessageInput).Return(mockContext.receiveMessageOutput, nil)
	mockContext.processor.EXPECT().ProcessEvent(*mockContext.receiveMessageOutput.Messages[0].Body).Return(errors.New("Process event failed"))
	mockContext.processor.EXPECT().ProcessEvent(*mockContext.receiveMessageOutput.Messages[1].Body).Return(nil)
	mockContext.sqsClient.EXPECT().DeleteMessage(mockContext.deleteMessageInput2).Return(nil, nil).Do(func(x interface{}) {
		pollCount++
		if pollCount == 1 {
			cancel()
		}
	})

	c.PollForEvents(ctx)
}

func TestPollForEventsFirstDeleteEventFails(t *testing.T) {
	mockContext := NewConsumerMockContext(t)
	defer mockContext.mockCtrl.Finish()

	mockContext.sqsClient.EXPECT().GetQueueUrl(gomock.Eq(mockContext.getQueueUrlInput)).Return(mockContext.getQueueUrlOutput, nil)

	c, err := NewConsumer(mockContext.sqsClient, mockContext.processor)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	pollCount := 0

	mockContext.sqsClient.EXPECT().ReceiveMessage(mockContext.receiveMessageInput).Return(mockContext.receiveMessageOutput, nil)
	mockContext.processor.EXPECT().ProcessEvent(*mockContext.receiveMessageOutput.Messages[0].Body).Return(nil)
	mockContext.sqsClient.EXPECT().DeleteMessage(mockContext.deleteMessageInput).Return(nil, errors.New("Delete event failed"))
	mockContext.processor.EXPECT().ProcessEvent(*mockContext.receiveMessageOutput.Messages[1].Body).Return(nil)
	mockContext.sqsClient.EXPECT().DeleteMessage(mockContext.deleteMessageInput2).Return(nil, nil).Do(func(x interface{}) {
		pollCount++
		if pollCount == 1 {
			cancel()
		}
	})

	c.PollForEvents(ctx)
}

func TestPollForEventsReceiveTwoMessages(t *testing.T) {
	mockContext := NewConsumerMockContext(t)
	defer mockContext.mockCtrl.Finish()

	mockContext.sqsClient.EXPECT().GetQueueUrl(gomock.Eq(mockContext.getQueueUrlInput)).Return(mockContext.getQueueUrlOutput, nil)

	c, err := NewConsumer(mockContext.sqsClient, mockContext.processor)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	pollCount := 0

	mockContext.sqsClient.EXPECT().ReceiveMessage(mockContext.receiveMessageInput).Return(mockContext.receiveMessageOutput, nil).Times(2)
	mockContext.processor.EXPECT().ProcessEvent(*mockContext.receiveMessageOutput.Messages[0].Body).Return(nil).Times(2)
	mockContext.sqsClient.EXPECT().DeleteMessage(mockContext.deleteMessageInput).Return(nil, nil).Times(2)
	mockContext.processor.EXPECT().ProcessEvent(*mockContext.receiveMessageOutput.Messages[1].Body).Return(nil).Times(2)
	mockContext.sqsClient.EXPECT().DeleteMessage(mockContext.deleteMessageInput2).Return(nil, nil).Times(2).Do(func(x interface{}) {
		pollCount++
		if pollCount == 2 {
			cancel()
		}
	})

	c.PollForEvents(ctx)
}
