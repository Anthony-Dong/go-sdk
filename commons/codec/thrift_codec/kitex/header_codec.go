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

// Copyright  https://github.com/cloudwego/kitex/blob/develop/pkg/remote/codec/header_codec.go

import (
	"fmt"
	"io"
)

/**
 *  TTHeader Protocol
 *  +-------------2Byte--------------|-------------2Byte-------------+
 *	+----------------------------------------------------------------+
 *	| 0|                          LENGTH                             |
 *	+----------------------------------------------------------------+
 *	| 0|       HEADER MAGIC          |            FLAGS              |
 *	+----------------------------------------------------------------+
 *	|                         SEQUENCE NUMBER                        |
 *	+----------------------------------------------------------------+
 *	| 0|     Header Size(/32)        | ...
 *	+---------------------------------
 *
 *	Header is of variable size:
 *	(and starts at offset 14)
 *
 *	+----------------------------------------------------------------+
 *	| PROTOCOL ID  |NUM TRANSFORMS . |TRANSFORM 0 ID (uint8)|
 *	+----------------------------------------------------------------+
 *	|  TRANSFORM 0 DATA ...
 *	+----------------------------------------------------------------+
 *	|         ...                              ...                   |
 *	+----------------------------------------------------------------+
 *	|        INFO 0 ID (uint8)      |       INFO 0  DATA ...
 *	+----------------------------------------------------------------+
 *	|         ...                              ...                   |
 *	+----------------------------------------------------------------+
 *	|                                                                |
 *	|                              PAYLOAD                           |
 *	|                                                                |
 *	+----------------------------------------------------------------+
 */

// Header keys
const (
	// Header Magics
	// 0 and 16th bits must be 0 to differentiate from framed & unframed
	TTHeaderMagic     uint32 = 0x10000000
	MeshHeaderMagic   uint32 = 0xFFAF0000
	MeshHeaderLenMask uint32 = 0x0000FFFF

	// HeaderMask        uint32 = 0xFFFF0000
	FlagsMask     uint32 = 0x0000FFFF
	MethodMask    uint32 = 0x41000000 // method first byte [A-Za-z_]
	MaxFrameSize  uint32 = 0x3FFFFFFF
	MaxHeaderSize uint32 = 65536
)

type HeaderFlags uint16

const (
	HeaderFlagsKey              string      = "HeaderFlags"
	HeaderFlagSupportOutOfOrder HeaderFlags = 0x01
	HeaderFlagDuplexReverse     HeaderFlags = 0x08
	HeaderFlagSASL              HeaderFlags = 0x10
)

const (
	TTHeaderMetaSize = 14
)

// ProtocolID is the wrapped protocol id used in THeader.
type ProtocolID uint8

// Supported ProtocolID values.
const (
	ProtocolIDThriftBinary  ProtocolID = 0x00
	ProtocolIDThriftCompact ProtocolID = 0x02
	ProtocolIDProtobufKitex ProtocolID = 0x03
	ProtocolIDDefault                  = ProtocolIDThriftBinary
)

type InfoIDType uint8 // uint8

const (
	InfoIDPadding     InfoIDType = 0
	InfoIDKeyValue    InfoIDType = 0x01
	InfoIDIntKeyValue InfoIDType = 0x10
	InfoIDACLToken    InfoIDType = 0x11
)

func readKVInfo(idx int, buf []byte, message *MetaInfo) error {
	intInfo := message.IntInfo
	strInfo := message.StrInfo
	for {
		infoID, err := Bytes2Uint8(buf, idx)
		idx++
		if err != nil {
			// this is the last field, read until there is no more padding
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		switch InfoIDType(infoID) {
		case InfoIDPadding:
			continue
		case InfoIDKeyValue:
			_, err := readStrKVInfo(&idx, buf, strInfo)
			if err != nil {
				return err
			}
		case InfoIDIntKeyValue:
			_, err := readIntKVInfo(&idx, buf, intInfo)
			if err != nil {
				return err
			}
		case InfoIDACLToken:
			err = skipACLToken(&idx, buf)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid infoIDType[%#x]", infoID)
		}
	}
	return nil
}

func readIntKVInfo(idx *int, buf []byte, info map[uint16]string) (has bool, err error) {
	kvSize, err := Bytes2Uint16(buf, *idx)
	*idx += 2
	if err != nil {
		return false, fmt.Errorf("error reading int kv info size: %s", err.Error())
	}
	if kvSize <= 0 {
		return false, nil
	}
	for i := uint16(0); i < kvSize; i++ {
		key, err := Bytes2Uint16(buf, *idx)
		*idx += 2
		if err != nil {
			return false, fmt.Errorf("error reading int kv info: %s", err.Error())
		}
		val, n, err := ReadString2BLen(buf, *idx)
		*idx += n
		if err != nil {
			return false, fmt.Errorf("error reading int kv info: %s", err.Error())
		}
		info[key] = val
	}
	return true, nil
}

func readStrKVInfo(idx *int, buf []byte, info map[string]string) (has bool, err error) {
	kvSize, err := Bytes2Uint16(buf, *idx)
	*idx += 2
	if err != nil {
		return false, fmt.Errorf("error reading str kv info size: %s", err.Error())
	}
	if kvSize <= 0 {
		return false, nil
	}
	for i := uint16(0); i < kvSize; i++ {
		key, n, err := ReadString2BLen(buf, *idx)
		*idx += n
		if err != nil {
			return false, fmt.Errorf("error reading str kv info: %s", err.Error())
		}
		val, n, err := ReadString2BLen(buf, *idx)
		*idx += n
		if err != nil {
			return false, fmt.Errorf("error reading str kv info: %s", err.Error())
		}
		info[key] = val
	}
	return true, nil
}

// skipACLToken SDK don't need acl token, just skip it
func skipACLToken(idx *int, buf []byte) error {
	_, n, err := ReadString2BLen(buf, *idx)
	*idx += n
	if err != nil {
		return fmt.Errorf("error reading acl token: %s", err.Error())
	}
	return nil
}
