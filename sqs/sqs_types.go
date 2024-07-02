package sqs

// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ReceiveMessage.html
type SqsReceiveMessageParams struct {
	AttributeNames              string `json:"AttributeNames,omitempty"`
	MaxNumberOfMessages         int64  `json:"MaxNumberOfMessages,omitempty"`
	MessageAttributeNames       string `json:"MessageAttributeNames,omitempty"`
	MessageSystemAttributeNames string `json:"MessageSystemAttributeNames,omitempty"`
	QueueUrl                    string `json:"QueueUrl"`
	ReceiveRequestAttemptId     string `json:"ReceiveRequestAttemptId,omitempty"`
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

type SqsMessage struct {
	Attributes             map[string]string           `json:"Attributes"`
	Body                   string                      `json:"Body"`
	MD5OfBody              string                      `json:"MD5OfBody"`
	MD5OfMessageAttributes string                      `json:"MD5OfMessageAttributes"`
	MessageAttributes      map[string]MessageAttribute `json:"MessageAttributes"`
	MessageId              string                      `json:"MessageId"`
	ReceiptHandle          string                      `json:"ReciptHandle"`
}

type SqsReceiveMessageResponse struct {
	Messages []SqsMessage `json:"Messages"`
}

// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_DeleteMessage.html
type SqsDeleteMessageParams struct {
	QueueUrl      string `json:"QueueUrl"`
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

type SqsSendMessageParams struct {
	DelaySeconds            int64                             `json:"DelaySeconds,omitempty"`
	MessageAttributes       map[string]MessageAttribute       `json:"MessageAttributes,omitempty"`
	MessageBody             string                            `json:"MessageBody"`
	MessageDeduplicationId  string                            `json:"MessageDeduplicationId,omitempty"`
	MessageGroupId          string                            `json:"MessageGroupId,omitempty"`
	MessageSystemAttributes map[string]MessageSystemAttribute `json:"MessageSystemAttributes,omitempty"`
	QueueUrl                string                            `json:"QueueUrl"`
}

type SqsSendMessageResponse struct {
	MD5OfMessageAttributes string `json:"MD5OfMessageAttributes"`
	MD5OfMessageBody       string `json:"MD5OfMessageBody"`
	MessageId              string `json:"MessageId"`
}
