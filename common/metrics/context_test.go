// Copyright (c) 2021 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package metrics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextTags(t *testing.T) {
	ctx := context.Background()
	assert.Empty(t, GetContextTags(ctx))

	tag1 := DomainTag("domain")
	ctx = TagContext(ctx, tag1)
	assert.Equal(t, []Tag{tag1}, GetContextTags(ctx))

	tag2 := ThriftTransportTag()
	ctx = TagContext(ctx, tag2)
	assert.Equal(t, []Tag{tag1, tag2}, GetContextTags(ctx))

	ctx1 := context.Background()
	ctx1 = TagContexts(ctx1, tag1, tag2)
	assert.Contains(t, GetContextTags(ctx1), tag1)
	assert.Contains(t, GetContextTags(ctx1), tag2)
}
