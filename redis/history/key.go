package history

import "fmt"

type CommonProxyStatus map[string]string

func GetPoolStatsKey(address string) string {
	return fmt.Sprintf("pool_stats_%s", address)
}

func GetSQLDigestKey(digest string) string {
	return fmt.Sprintf("sql_digest_%s", digest)
}

func GetEndpointStatsKey(address string) string {
	return fmt.Sprintf("endpoint_stats_%s", address)
}

func GetEndpointLogKey(address string) string {
	return fmt.Sprintf("endpoint_log_%s", address)
}

func GetRouteRuleHitKey(address string) string {
	return fmt.Sprintf("route_rule_hit_%s", address)
}

func GetSQLStatDigestKey(address, dbname, digest string) string {
	if dbname == "" {
		dbname = "default"
	}
	return fmt.Sprintf("sql_stat_%s_%s_digest_%s", address, dbname, digest)
}

func GetProxyDigestKey(address string) string {
	return fmt.Sprintf("sql_stat_%s", address)
}
