/*
   Copyright 2014 CoreOS, Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package httptypes

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMemberUnmarshal(t *testing.T) {
	tests := []struct {
		body       []byte
		wantMember Member
		wantError  bool
	}{
		// no URLs, just check ID & Name
		{
			body:       []byte(`{"id": 1, "name": "dungarees"}`),
			wantMember: Member{ID: 1, Name: "dungarees", PeerURLs: nil, ClientURLs: nil},
		},

		// both client and peer URLs
		{
			body: []byte(`{"peerURLs": ["http://127.0.0.1:4001"], "clientURLs": ["http://127.0.0.1:4001"]}`),
			wantMember: Member{
				PeerURLs: []string{
					"http://127.0.0.1:4001",
				},
				ClientURLs: []string{
					"http://127.0.0.1:4001",
				},
			},
		},

		// multiple peer URLs
		{
			body: []byte(`{"peerURLs": ["http://127.0.0.1:4001", "https://example.com"]}`),
			wantMember: Member{
				PeerURLs: []string{
					"http://127.0.0.1:4001",
					"https://example.com",
				},
				ClientURLs: nil,
			},
		},

		// multiple client URLs
		{
			body: []byte(`{"clientURLs": ["http://127.0.0.1:4001", "https://example.com"]}`),
			wantMember: Member{
				PeerURLs: nil,
				ClientURLs: []string{
					"http://127.0.0.1:4001",
					"https://example.com",
				},
			},
		},

		// invalid JSON
		{
			body:      []byte(`{"peerU`),
			wantError: true,
		},
	}

	for i, tt := range tests {
		got := Member{}
		err := json.Unmarshal(tt.body, &got)
		if tt.wantError != (err != nil) {
			t.Errorf("#%d: want error %t, got %v", i, tt.wantError, err)
			continue
		}

		if !reflect.DeepEqual(tt.wantMember, got) {
			t.Errorf("#%d: incorrect output: want=%#v, got=%#v", i, tt.wantMember, got)
		}
	}
}

func TestMemberCollectionUnmarshal(t *testing.T) {
	body := []byte(`{"members":[{"id":176869799018424574,"peerURLs":["http://127.0.0.1:7003"],"name":"node3","clientURLs":["http://127.0.0.1:4003"]},{"id":297577273835923749,"peerURLs":["http://127.0.0.1:2380","http://127.0.0.1:7001"],"name":"node1","clientURLs":["http://127.0.0.1:2379","http://127.0.0.1:4001"]},{"id":10666918107976480891,"peerURLs":["http://127.0.0.1:7002"],"name":"node2","clientURLs":["http://127.0.0.1:4002"]}]}`)

	want := MemberCollection{
		Members: []Member{
			{
				ID:   176869799018424574,
				Name: "node3",
				PeerURLs: []string{
					"http://127.0.0.1:7003",
				},
				ClientURLs: []string{
					"http://127.0.0.1:4003",
				},
			},
			{
				ID:   297577273835923749,
				Name: "node1",
				PeerURLs: []string{
					"http://127.0.0.1:2380",
					"http://127.0.0.1:7001",
				},
				ClientURLs: []string{
					"http://127.0.0.1:2379",
					"http://127.0.0.1:4001",
				},
			},
			{
				ID:   10666918107976480891,
				Name: "node2",
				PeerURLs: []string{
					"http://127.0.0.1:7002",
				},
				ClientURLs: []string{
					"http://127.0.0.1:4002",
				},
			},
		},
	}

	got := MemberCollection{}
	err := json.Unmarshal(body, &got)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Incorrect output: want=%#v, got=%#v", want, got)
	}
}
