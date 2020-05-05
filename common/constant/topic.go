package constant

// topic 消息类型
const (
	TOPIC_USER_LOGIN   string = "topic.user.login"
	TOPIC_USER_CREATED string = "topic.user.created"

	// G2L: gate(RPC) -> ums(MQ: TOPIC_L2A_PREFIX) -> chat
	TOPIC_L2A_PREFIX      string = "topic.L2A:%d"      // %d为appid
	// L2G: chat(RPC) -> ums(MQ: TOPIC_L2G_PREFIX) -> web/ws
	TOPIC_L2G_PREFIX      string = "topic.L2G:%d"      // %d为gateid
	TOPIC_L2G_PUSH_PREFIX string = "topic.L2G.push:%d" // %d为gateid

)
