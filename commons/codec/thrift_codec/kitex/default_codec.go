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

// Copyright https://github.com/cloudwego/kitex/blob/develop/pkg/remote/codec/default_codec.go

import (
	"encoding/binary"
	"fmt"
)

// The byte count of 32 and 16 integer values.
const (
	Size32 = 4
	Size16 = 2
)

const (
	// ThriftV1Magic is the magic code for thrift.VERSION_1
	ThriftV1Magic = 0x80010000
	// ProtobufV1Magic is the magic code for kitex protobuf
	ProtobufV1Magic = 0x90010000

	// MagicMask is bit mask for checking version.
	MagicMask = 0xffff0000
)

// NewDefaultCodec creates the default protocol sniffing codec supporting thrift and protobuf.

/**
 * +------------------------------------------------------------+
 * |                  4Byte                 |       2Byte       |
 * +------------------------------------------------------------+
 * |   			     Length			    	|   HEADER MAGIC    |
 * +------------------------------------------------------------+
 */
func isTTHeader(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[Size32:])&MagicMask == TTHeaderMagic
}

/**
 * +----------------------------------------+
 * |       2Byte        |       2Byte       |
 * +----------------------------------------+
 * |    HEADER MAGIC    |   HEADER SIZE     |
 * +----------------------------------------+
 */
func isMeshHeader(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[:Size32])&MagicMask == MeshHeaderMagic
}

/**
 * Kitex protobuf has framed field
 * +------------------------------------------------------------+
 * |                  4Byte                 |       2Byte       |
 * +------------------------------------------------------------+
 * |   			     Length			    	|   HEADER MAGIC    |
 * +------------------------------------------------------------+
 */
func isProtobufKitex(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[Size32:])&MagicMask == ProtobufV1Magic
}

/**
 * +-------------------+
 * |       2Byte       |
 * +-------------------+
 * |   HEADER MAGIC    |
 * +-------------------
 */
func isThriftBinary(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[:Size32])&MagicMask == ThriftV1Magic
}

/**
 * +------------------------------------------------------------+
 * |                  4Byte                 |       2Byte       |
 * +------------------------------------------------------------+
 * |   			     Length			    	|   HEADER MAGIC    |
 * +------------------------------------------------------------+
 */
func isThriftFramedBinary(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[Size32:])&MagicMask == ThriftV1Magic
}

// protoID just for ttheader
func CheckProtocolID(protoID uint8) error {
	switch protoID {
	case uint8(ProtocolIDProtobufKitex):
		// rpcCfg := internal.AsMutableRPCConfig(message.RPCInfo().Config())
		// rpcCfg.SetCodecType(kitex.TTHeaderProtobufKitex)
	case uint8(ProtocolIDThriftBinary):
	default:
		return fmt.Errorf("unsupport ProtocolID[%d]", protoID)
	}
	return nil
}
