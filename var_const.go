package lfc

const (
	redis_key_prefix        = "fc" //redis key prefix
	redis_key_shadow_suffix = "_shadow"
	redis_key_split         = ":" //redis key delimiter
	redis_key_inner_split   = "#" //redis key delimiter : fc:1:uid#device_id:86400:[uid's value]#[device's value]
)
