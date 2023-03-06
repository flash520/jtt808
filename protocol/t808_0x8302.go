package protocol

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

// 提问下发
type T808_0x8302 struct {
	// 标志
	Flag byte
	// 问题
	Question string
	// 候选答案列表
	CandidateAnswers []T808_0x8302_Answer
}

// 问题答案
type T808_0x8302_Answer struct {
	// 答案 ID
	ID byte
	// 答案内容
	Content string
}

func (entity *T808_0x8302) MsgID() MsgID {
	return MsgT808_0x8302
}

func (entity *T808_0x8302) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入标志位
	writer.WriteByte(entity.Flag)

	// 写入问题长度
	reader := bytes.NewReader([]byte(entity.Question))
	question, err := ioutil.ReadAll(
		transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder()))
	if err != nil {
		return nil, err
	}
	writer.WriteByte(byte(len(question)))

	// 写入问题内容
	writer.Write(question)

	// 写入候选答案
	for _, answer := range entity.CandidateAnswers {
		// 写入答案ID
		writer.WriteByte(answer.ID)

		// 写入答案长度
		reader := bytes.NewReader([]byte(answer.Content))
		content, err := ioutil.ReadAll(
			transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder()))
		if err != nil {
			return nil, err
		}
		writer.WriteByte(byte(len(content)))

		// 写入答案内容
		writer.Write(content)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8302) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取标志位
	var err error
	entity.Flag, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取问题长度
	size, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取问题内容
	entity.Question, err = reader.ReadString(int(size))
	if err != nil {
		return 0, err
	}

	// 读取候选答案
	for {
		if reader.Len() == 0 {
			break
		}

		var answer T808_0x8302_Answer

		// 读取答案ID
		answer.ID, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取内容长度
		size, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取事件内容
		answer.Content, err = reader.ReadString(int(size))
		if err != nil {
			return 0, err
		}
		entity.CandidateAnswers = append(entity.CandidateAnswers, answer)
	}

	return len(data) - reader.Len(), nil
}
