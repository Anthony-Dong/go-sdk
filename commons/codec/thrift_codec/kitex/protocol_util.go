/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kitex

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
)

func NewMetaInfo() *MetaInfo {
	return &MetaInfo{
		IntInfo: map[uint16]string{},
		StrInfo: map[string]string{},
	}
}

type MetaInfo struct {
	IntInfo map[uint16]string
	StrInfo map[string]string
}

var (
	_intKeyMap = map[uint16]string{
		MeshVersion:     "MeshVersion",
		TransportType:   "TransportType",
		LogID:           "LogID",
		FromService:     "FromService",
		FromCluster:     "FromCluster",
		FromIDC:         "FromIDC",
		ToService:       "ToService",
		ToCluster:       "ToCluster",
		ToIDC:           "ToIDC",
		ToMethod:        "ToMethod",
		Env:             "Env",
		DestAddress:     "DestAddress",
		RPCTimeout:      "RPCTimeout",
		ReadTimeout:     "ReadTimeout",
		RingHashKey:     "RingHashKey",
		DDPTag:          "DDPTag",
		WithMeshHeader:  "WithMeshHeader",
		ConnectTimeout:  "ConnectTimeout",
		SpanContext:     "SpanContext",
		ShortConnection: "ShortConnection",
		FromMethod:      "FromMethod",
		StressTag:       "StressTag",
		MsgType:         "MsgType",
		HTTPContentType: "HTTPContentType",
		RawRingHashKey:  "RawRingHashKey",
		LBType:          "LBType",
	}
	_msgType = map[string]string{
		"0": "INVALID_TMESSAGE_TYPE", // thrift.INVALID_TMESSAGE_TYPE,
		"1": "CALL",                  // thrift.CALL,
		"2": "REPLY",                 // thrift.REPLY,
		"3": "EXCEPTION",             //thrift.EXCEPTION,
		"4": "ONEWAY",                // thrift.ONEWAY,
	}
)

func getMsgType(v string) string {
	if r, isExist := _msgType[v]; isExist {
		return r
	}
	return v
}
func getIntKey(k uint16) string {
	if sk, isExist := _intKeyMap[k]; isExist {
		return sk
	}
	return strconv.Itoa(int(k))
}

func (m MetaInfo) MarshalJSON() (text []byte, err error) {
	return []byte(m.String()), nil
}

func (m MetaInfo) String() string {
	result := make(map[string]map[string]string, 0)
	if len(m.StrInfo) != 0 {
		result["StrInfo"] = m.StrInfo
	}
	intInfo := make(map[string]string, len(m.IntInfo))
	for k, v := range m.IntInfo {
		if k == MsgType {
			v = getMsgType(v)
		}
		if k == WithMeshHeader {
			if v == MeshOrTTHeaderProtocol {
				v = "MeshOrTTHeaderProtocol"
			}
		}
		intInfo[getIntKey(k)] = v
	}
	if len(intInfo) != 0 {
		result["IntInfo"] = intInfo
	}
	marshal, _ := json.Marshal(result)
	return string(marshal)
}

func IsTTHeader(reader *bufio.Reader) bool {
	flag, err := reader.Peek(Size32 * 2)
	if err != nil {
		return false
	}
	return isTTHeader(flag)
}

func ReadTTHeader(reader *bufio.Reader, metaInfo *MetaInfo) (headerSize int, err error) {
	headerMeta, err := reader.Peek(TTHeaderMetaSize)
	if err != nil {
		return 0, err
	}
	headerInfoSize := binary.BigEndian.Uint16(headerMeta[Size32*3:TTHeaderMetaSize]) * 4
	if uint32(headerInfoSize) > MaxHeaderSize || headerInfoSize < 2 {
		return 0, err
	}
	// append headerInfo
	headerInfo, err := reader.Peek(TTHeaderMetaSize + int(headerInfoSize))
	if err != nil {
		return 0, fmt.Errorf("invalid header length[%d]", headerInfoSize)
	}
	headerInfo = headerInfo[TTHeaderMetaSize:]
	if err := CheckProtocolID(headerInfo[0]); err != nil {
		return 0, err
	}
	hdIdx := 2
	transformIDNum := int(headerInfo[1])
	if int(headerInfoSize)-hdIdx < transformIDNum {
		return 0, fmt.Errorf("need read %d transformIDs, but not enough", transformIDNum)
	}
	transformIDs := make([]uint8, transformIDNum)
	for i := 0; i < transformIDNum; i++ {
		transformIDs[i] = headerInfo[hdIdx]
		hdIdx++
	}
	if err := readMetaInfo(hdIdx, headerInfo, metaInfo); err != nil {
		return 0, fmt.Errorf("ttHeader read kv info failed, %s", err)
	}
	return TTHeaderMetaSize + int(headerInfoSize), nil
}

func readMetaInfo(idx int, buf []byte, message *MetaInfo) error {
	if message == nil {
		message = NewMetaInfo()
	}
	if message.IntInfo == nil {
		message.IntInfo = map[uint16]string{}
	}
	if message.StrInfo == nil {
		message.StrInfo = map[string]string{}
	}
	return readKVInfo(idx, buf, message)
}

func IsMeshHeader(reader *bufio.Reader) bool {
	flag, err := reader.Peek(2 * Size32)
	if err != nil {
		return false
	}
	return isMeshHeader(flag)
}

func ReadMeshHeader(reader *bufio.Reader, metaInfo *MetaInfo) (size int, err error) {
	headerMeta, err := reader.Peek(Size32)
	if err != nil {
		return 0, fmt.Errorf("meshHeader read header meta failed: %v", err)
	}
	headerLen := binary.BigEndian.Uint16(headerMeta[Size16:])
	headerInfo, err := reader.Peek(Size32 + int(headerLen))
	if err != nil {
		return 0, fmt.Errorf("meshHeader read header buf failed: %s", err.Error())
	}
	headerInfo = headerInfo[Size32:]
	idx := 0
	if metaInfo == nil {
		metaInfo = NewMetaInfo()
	}
	if _, err = readStrKVInfo(&idx, headerInfo, metaInfo.StrInfo); err != nil {
		return 0, fmt.Errorf("meshHeader read kv info failed, %s", err.Error())
	}
	return Size32 + int(headerLen), nil
}
