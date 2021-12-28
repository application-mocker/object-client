package object_client

import (
	"fmt"
	"testing"
)

func TestNewObjectClient(t *testing.T) {
	type args struct {
		host      string
		baseScope string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test common", args: args{host: "http://127.0.0.1:3000", baseScope: "test"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewObjectClient(tt.args.host, tt.args.baseScope)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewObjectClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("Return is nil")
			}
		})
	}
}

func TestObjectClient_InsertOne(t *testing.T) {
	client, err := NewObjectClient("http://127.0.0.1:3000", "test_scope")

	if err != nil {
		fmt.Println(err)
	}
	obj := struct {
		Name string `json:"name"`
	}{
		Name: "Test",
	}

	id, err := client.InsertOne(obj)
	fmt.Println(id)

}

func TestObjectClient_GetById(t *testing.T) {
	client, err := NewObjectClient("http://127.0.0.1:3000", "test_scope")

	if err != nil {
		fmt.Println(err)
	}
	obj := &struct {
		BaseNode
		Data struct {
			Name string `json:"name"`
		} `json:"data"`
	}{}

	id, err := client.GetById("bfd87b46-67aa-11ec-8000-acde48001122", obj)
	fmt.Println(id)
	fmt.Println(obj)
}

func TestObjectClient_GetByIdWithoutBaseStruct(t *testing.T) {
	client, err := NewObjectClient("http://127.0.0.1:3000", "test_scope")

	if err != nil {
		fmt.Println(err)
	}
	obj := &struct {
		Name string `json:"name"`
	}{}

	id, err := client.GetByIdWithoutBaseStruct("bfd87b46-67aa-11ec-8000-acde48001122", obj)
	fmt.Println(id)
	fmt.Println(obj)
}

func TestObjectClient_ListAllValue(t *testing.T) {
	client, err := NewObjectClient("http://127.0.0.1:3000", "test_scope")

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.ListAllValue()
	if res != nil {
		for _, item := range res {
			fmt.Println(item.Id, item.DataValue)
		}
	}
	fmt.Println(err)
}

func TestObjectClient_DeleteById(t *testing.T) {
	client, err := NewObjectClient("http://127.0.0.1:3000", "test_scope")

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.DeleteById("296bf544-67a4-11ec-8000-acde48001122")
	if res != nil {
		fmt.Println(*res)
	}

	fmt.Println(err)
}

func TestObjectClient_UpdateByIdWithoutBaseStruct(t *testing.T) {
	client, err := NewObjectClient("http://127.0.0.1:3000", "test_scope")

	if err != nil {
		fmt.Println(err)
	}
	obj := &struct {
		Name string `json:"name"`
		Set  string `json:"TEests"`
	}{
		Name: "Test",
		Set:  "tsakjfdaslfjasfj asf;alfj asdjf ;aljf ;asf ;lasjf l;",
	}

	res, err := client.UpdateByIdWithoutBaseStruct("bfd87b46-67aa-11ec-8000-acde48001122", obj)
	fmt.Println(res)
	fmt.Println(obj)
}
