package design // The convention consists of naming the design
// package "design"
import (
	. "github.com/goadesign/goa/design" // Use . imports to enable the DSL
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("ddgo", func() { // API defines the microservice endpoint and
	Title("Make monitoring setting on DD") // other global properties. There should be one
	Description("")                        // and exactly one API definition appearing in
	Scheme("http")                         // the design.
	Host("localhost:8080")
})

var _ = Resource("monitor", func() { // Resources group related API endpoints
	BasePath("/monitors")      // together. They map to REST resources for REST
	DefaultMedia(MonitorMedia) // services.

	Action("getAllMonitorDetails", func() { // Actions define a single API endpoint together
		Description("Get all monitoring setting") // with its path, parameters (both path
		// Routing(GET("/:bottleID"))                  // parameters and querystring values) and payload
		Params(func() { // (shape of the request body).
			// Param("group_states", Integer, "Bottle ID")
		})
		Response(OK)       // Responses define the shape and status code
		Response(NotFound) // of HTTP responses.
	})
})

// UserMedia defines the media type used to render users.
var MonitorMedia = MediaType("application/vnd.goa.ddgo.monitor+json", func() {
	Description("A setting of monitor")
	Attributes(func() { // Attributes define the media type shape.
		Attribute("id", Integer, "Unique monitori ID")
		Attribute("message", String, "Alert message which includes notified 	target")
		Attribute("name", String, "Name of wine")
		Required("id", "message", "name")
	})
	View("default", func() { // View defines a rendering of the media type.
		Attribute("id")      // Media types may have multiple views and must
		Attribute("message") // have a "default" view.
		Attribute("name")
	})
})

/*
- 作成し終わった後、過不足をチェックする
-- Integrationを取得するためのAPIは存在しない
  - How to Validate
    - Monitorを全部取ってきて、Integrationが既存のものかを確認する -> NG 複数回実行されたら詰む

ま、いっか。テンプレからスケルトンを生成して、一覧を渡すからちゃんとIntegration設定とテストしなよ！ってことにしーようっと。テストされてないよ通知Channel作れば問題ないっしょ

1. User APIでIntegrationの有無を確認する
2. 無いものがあれば、res

*/
