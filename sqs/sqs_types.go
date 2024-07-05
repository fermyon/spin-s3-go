package sqs

// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ReceiveMessage.html
type ReceiveMessageParams struct {
	AttributeNames              string `json:"AttributeNames,omitempty"`
	MaxNumberOfMessages         int64  `json:"MaxNumberOfMessages,omitempty"`
	MessageAttributeNames       string `json:"MessageAttributeNames,omitempty"`
	MessageSystemAttributeNames string `json:"MessageSystemAttributeNames,omitempty"`
	QueueURL                    string `json:"QueueUrl"`
	ReceiveRequestAttemptID     string `json:"ReceiveRequestAttemptId,omitempty"`
	VisibilityTimeout           int64  `json:"VisibilityTimeout,omitempty"`
}

type MessageAttribute struct {
	// TODO: @asteurer Validate ChatGPT suggestion
	BinaryListValues [][]byte `json:"BinaryListValues"`
	BinaryValue      []byte   `json:"BinaryValue"`
	DataType         string   `json:"DataType"`
	StringListValues []string `json:"StringListValues"`
	StringValue      string   `json:"StringValue"`
}

type Message struct {
	Attributes             map[string]string           `json:"Attributes"`
	Body                   string                      `json:"Body"`
	MD5OfBody              string                      `json:"MD5OfBody"`
	MD5OfMessageAttributes string                      `json:"MD5OfMessageAttributes"`
	MessageAttributes      map[string]MessageAttribute `json:"MessageAttributes"`
	MessageID              string                      `json:"MessageId"`
	ReceiptHandle          string                      `json:"ReciptHandle"`
}

type ReceiveMessageResponse struct {
	Messages []Message `json:"Messages"`
}

// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_DeleteMessage.html
type DeleteMessageParams struct {
	QueueURL      string `json:"QueueUrl"`
	ReceiptHandle string `json:"ReceiptHandle"`
}

// There is not a JSON body response for a DeleteMessage request

// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_SendMessage.html
type MessageSystemAttribute struct {
	BinaryListValues [][]byte `json:"BinaryListValues"`
	BinaryValue      []byte   `json:"BinaryValue"`
	DataType         string   `json:"DataType"`
	StringListValues []string `json:"StringListValues"`
	StringValue      string   `json:"StringValue"`
}

type SendMessageParams struct {
	DelaySeconds            int64                             `json:"DelaySeconds,omitempty"`
	MessageAttributes       map[string]MessageAttribute       `json:"MessageAttributes,omitempty"`
	MessageBody             string                            `json:"MessageBody"`
	MessageDeduplicationID  string                            `json:"MessageDeduplicationId,omitempty"`
	MessageGroupID          string                            `json:"MessageGroupId,omitempty"`
	MessageSystemAttributes map[string]MessageSystemAttribute `json:"MessageSystemAttributes,omitempty"`
	QueueURL                string                            `json:"QueueUrl"`
}

type SendMessageResponse struct {
	MD5OfMessageAttributes string `json:"MD5OfMessageAttributes,omitempty"`
	MD5OfMessageBody       string `json:"MD5OfMessageBody"`
	MessageID              string `json:"MessageId"`
}

// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_CreateQueue.html
type CreateQueueParams struct {
	Attributes map[string]string `json:"Attributes,omitempty"`
	QueueName  string            `json:"QueueName"`
	Tags       map[string]string `json:"Tags,omitempty"`
}

type CreateQueueResponse struct {
	QueueURL string `json:"QueueUrl"`
}
