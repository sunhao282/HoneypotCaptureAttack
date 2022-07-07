package redis

import (
	"bufio"
	"bytes"
	"fmt"
	"honey/src/core/pool"
	"honey/src/util/log"
	"honey/src/util/try"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var kvData map[string]string

func Start(addr string) {
	kvData = make(map[string]string)
	netListen, _ := net.Listen("tcp", addr)
	wg, poolX := pool.New(10)

	defer poolX.Release()
	defer netListen.Close()
	for {
		wg.Add(1)
		poolX.Submit(func() {
			time.Sleep(time.Second * 2)

			conn, err := netListen.Accept()

			if err != nil {
				log.Pr("Redis", "127.0.0.1", "Redis fail", err)
			}

			arr := strings.Split(conn.RemoteAddr().String(), ":")

			log.Pr("Redis", arr[0], "connect")

			go handleConnection(conn)

			wg.Done()
		})
	}
}

func handleConnection(conn net.Conn) {

	for {
		str := parseRESP(conn)
		arr := strings.Split(conn.RemoteAddr().String(), ":")
		switch value := str.(type) {
		case string:
			log.Pr("Redis", arr[0], "cmd", str)
			if len(value) == 0 {
				goto end
			}

			conn.Write([]byte(value))

		case []string:
			log.Pr("Redis", arr[0], "cmd", str)
			if strings.EqualFold(value[0], "set") {
				// 模拟 redis set

				try.Try(func() {
					key := string(value[1])
					val := string(value[2])
					kvData[key] = val

				}).Catch(func() {
					// 取不到 key 会异常
				})

				conn.Write([]byte("+OK\r\n"))
			} else if strings.EqualFold(value[0], "get") {
				try.Try(func() {
					// 模拟 redis get
					key := string(value[1])
					val := string(kvData[key])

					valLen := strconv.Itoa(len(val))
					if len(val) == 0 {
						str := "$" + "-1" + "\r\n"
						conn.Write([]byte(str))
					} else {
						str := "$" + valLen + "\r\n" + val + "\r\n"

						conn.Write([]byte(str))
					}
				}).Catch(func() {
					conn.Write([]byte("+OK\r\n"))
				})
			} else if strings.EqualFold(value[0], "info") {
				try.Try(func() {
					str := "$2748" + "\r\n" + "# Server" + "\r\n" + "redis_version:4.0.11" + "\r\n" + "redis_git_sha1:00000000" + "\r\n" + "redis_git_dirty:0" + "\r\n" + "redis_build_id:dc87299ecff67acc" + "\r\n" + "redis_mode:standalone" + "\r\n" + "os:Linux 3.10.0-957.el7.x86_64 x86_64" + "\r\n" + "arch_bits:64" + "\r\n" + "multiplexing_api:epoll" + "\r\n" + "atomicvar_api:atomic-builtin" + "\r\n" + "gcc_version:4.8.5" + "\r\n" + "process_id:37910" + "\r\n" + "run_id:739f37b5ed0173c1bb63df5752be21f0c04b2370" + "\r\n" + "tcp_port:6379" + "\r\n" + "uptime_in_seconds:262" + "\r\n" + "uptime_in_days:0" + "\r\n" + "hz:10" + "\r\n" + "lru_clock:368408" + "\r\n" + "executable:/usr/local/redis/bin/./redis-server" + "\r\n" + "config_file:/usr/local/redis/redis.conf" + "\r\n" + "\r\n" + "# Clients" + "\r\n" + "connected_clients:1" + "\r\n" + "client_longest_output_list:0" + "\r\n" + "client_biggest_input_buf:0" + "\r\n" + "blocked_clients:0" + "\r\n" + "\r\n" + "# Memory" + "\r\n" + "used_memory:849424" + "\r\n" + "used_memory_human:829.52K" + "\r\n" + "used_memory_rss:1810432" + "\r\n" + "used_memory_rss_human:1.73M" + "\r\n" + "used_memory_peak:849424" + "\r\n" + "used_memory_peak_human:829.52K" + "\r\n" + "used_memory_peak_perc:105.20%" + "\r\n" + "used_memory_overhead:836190" + "\r\n" + "used_memory_startup:786560" + "\r\n" + "used_memory_dataset:13234" + "\r\n" + "used_memory_dataset_perc:21.05%" + "\r\n" + "total_system_memory:1019797504" + "\r\n" + "total_system_memory_human:972.55M" + "\r\n" + "used_memory_lua:37888" + "\r\n" + "used_memory_lua_human:37.00K" + "\r\n" + "maxmemory:0" + "\r\n" + "maxmemory_human:0B" + "\r\n" + "maxmemory_policy:noeviction" + "\r\n" + "mem_fragmentation_ratio:2.13" + "\r\n"
					str2 := "mem_allocator:jemalloc-4.0.3" + "\r\n" + "active_defrag_running:0" + "\r\n" + "lazyfree_pending_objects:0" + "\r\n" + "\r\n" + "# Persistence" + "\r\n" + "loading:0" + "\r\n" + "rdb_changes_since_last_save:0" + "\r\n" + "rdb_bgsave_in_progress:0" + "\r\n" + "rdb_last_save_time:1577426450" + "\r\n" + "rdb_last_bgsave_status:ok" + "\r\n" + "rdb_last_bgsave_time_sec:-1" + "\r\n" + "rdb_current_bgsave_time_sec:-1" + "\r\n" + "rdb_last_cow_size:0" + "\r\n" + "aof_enabled:0" + "\r\n" + "aof_rewrite_in_progress:0" + "\r\n" + "aof_rewrite_scheduled:0" + "\r\n" + "aof_last_rewrite_time_sec:-1" + "\r\n" + "aof_current_rewrite_time_sec:-1" + "\r\n" + "aof_last_bgrewrite_status:ok" + "\r\n" + "aof_last_write_status:ok" + "\r\n" + "aof_last_cow_size:0" + "\r\n" + "\r\n" + "# Stats" + "\r\n" + "total_connections_received:1" + "\r\n" + "total_commands_processed:0" + "\r\n" + "instantaneous_ops_per_sec:0" + "\r\n" + "total_net_input_bytes:14" + "\r\n" + "total_net_output_bytes:0" + "\r\n" + "instantaneous_input_kbps:0.00" + "\r\n" + "instantaneous_output_kbps:0.00" + "\r\n" + "rejected_connections:0" + "\r\n" + "sync_full:0" + "\r\n" + "sync_partial_ok:0" + "\r\n" + "sync_partial_err:0" + "\r\n" + "expired_keys:0" + "\r\n" + "expired_stale_perc:0.00" + "\r\n" + "expired_time_cap_reached_count:0" + "\r\n" + "evicted_keys:0" + "\r\n" + "keyspace_hits:0" + "\r\n" + "keyspace_misses:0" + "\r\n" + "pubsub_channels:0" + "\r\n" + "pubsub_patterns:0" + "\r\n" + "latest_fork_usec:0" + "\r\n" + "migrate_cached_sockets:0" + "\r\n" + "slave_expires_tracked_keys:0" + "\r\n" + "active_defrag_hits:0" + "\r\n" + "active_defrag_misses:0" + "\r\n" + "active_defrag_key_hits:0" + "\r\n" + "active_defrag_key_misses:0" + "\r\n" + "\r\n" + "# Replication" + "\r\n" + "role:master" + "\r\n" + "connected_slaves:0" + "\r\n" + "master_replid:f4de5c09dbecc9962610e30d2d087932d140716e" + "\r\n" + "master_replid2:0000000000000000000000000000000000000000" + "\r\n" + "master_repl_offset:0" + "\r\n" + "second_repl_offset:-1" + "\r\n" + "repl_backlog_active:0" + "\r\n" + "repl_backlog_size:1048576" + "\r\n" + "repl_backlog_first_byte_offset:0" + "\r\n" + "repl_backlog_histlen:0" + "\r\n" + "\r\n" + "# CPU" + "\r\n" + "used_cpu_sys:0.27" + "\r\n" + "used_cpu_user:0.05" + "\r\n" + "used_cpu_sys_children:0.00" + "\r\n" + "used_cpu_user_children:0.00" + "\r\n" + "\r\n" + "# Cluster" + "\r\n" + "cluster_enabled:0" + "\r\n" + "\r\n" + "# Keyspace" + "\r\n" + "\r\n"
					str3 := str + str2
					fmt.Println(len(str3))
					conn.Write([]byte(str3))

				}).Catch(func() {
					// 取不到 key 会异常
				})

			} else if strings.EqualFold(value[0], "ping") {
				try.Try(func() {

					// fmt.Printf("redis  %s \n",value[0])
					str := "+PONG" + "\r\n"
					conn.Write([]byte(str))

				}).Catch(func() {
					// 取不到 key 会异常
				})

			} else if strings.EqualFold(value[0], "scan") {
				try.Try(func() {

					// fmt.Printf("redis  %s \n",value[0])
					str := "*2" + "\r\n" + "$1" + "\r\n" + "0" + "\r\n" + "*0" + "\r\n"
					conn.Write([]byte(str))

				}).Catch(func() {
					// 取不到 key 会异常
				})

			} else if strings.EqualFold(value[0], "client") {
				try.Try(func() {
					str2 := strip(conn.RemoteAddr().String(), "/ \n\r\t")

					str := "id=15 addr=" + str2 + " fd=8 name= age=426 idle=0 flags=N db=0 sub=0 psub=0 multi=-1 qbuf=0 qbuf-free=32768 obl=0 oll=0 omem=0 events=r cmd=client" + "\n\r\n"
					str_len := len(str) - 2
					str3 := "$" + strconv.Itoa(str_len) + "\r\n" + str
					//fmt.Println(str3)
					conn.Write([]byte(str3))

				}).Catch(func() {
					// 取不到 key 会异常
				})

			} else if strings.EqualFold(value[0], "del") {
				try.Try(func() {

					// fmt.Printf("redis  %s \n",value[0])
					str := ":1" + "\r\n"
					conn.Write([]byte(str))

				}).Catch(func() {
					// 取不到 key 会异常
				})

			} else if strings.EqualFold(value[0], "slaveof") {

				//	 fmt.Printf("redis  %s \n",value[0])
				//str := ":1"+"\r\

				conn.Write([]byte("+OK\r\n"))

				go handlerRedisConn(value[1], value[2])

			} else {
				try.Try(func() {

				}).Catch(func() {

				})

				conn.Write([]byte("+OK\r\n"))
			}
			break
		default:

		}
	}
end:
	conn.Close()
}
func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()

	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
	}

	return result.Bytes(), nil
}
func handlerRedisConn(ip string, port string) {
	//fmt.Println("redis ", id)
	str := ip + ":" + port
	redis_conn, err := net.Dial("tcp", str)
	if err != nil {
		log.Pr("Redis", ip, "Redis conn fail", str)
		return
	}
	redis_str := "PING" + "\r\n"
	redis_conn.Write([]byte(redis_str))
	strone := parseRedisConn(redis_conn)
	log.Pr("Redis_conn", str, "cmd", strone)
	time.Sleep(1 * time.Second)
	redis_str_one := "REPLCONF listening-port 6379" + "\r\n"
	redis_conn.Write([]byte(redis_str_one))
	strtwo := parseRedisConn(redis_conn)
	log.Pr("Redis_conn", str, "cmd", strtwo)
	time.Sleep(1 * time.Second)

	redis_str_two := "REPLCONF capa eof capa psync2" + "\r\n"
	redis_conn.Write([]byte(redis_str_two))
	str_thr := parseRedisConn(redis_conn)
	log.Pr("Redis_conn", str, "cmd", str_thr)
	time.Sleep(1 * time.Second)
	redis_str_three := "PSYNC a1fb3404be940fe378b00cdf5c4be33f5a4cebdc 1" + "\r\n"
	redis_conn.Write([]byte(redis_str_three))
	//str_for:=parseRedisConn( redis_conn)
	//log.Pr("Redis_conn", str, "cmd", str_for)
	result, err := readFully(redis_conn)
	checkError(err)
	filename := "redis" + time.Now().Format("2006-01-02-15-04-05")
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = f.Write(result)
	if err != nil {
		fmt.Println(err.Error())
	}
	redis_conn.Close()
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
func strip(s_ string, chars_ string) string {
	s, chars := []rune(s_), []rune(chars_)
	length := len(s)
	max := len(s) - 1
	l, r := true, true //标记当左端或者右端找到正常字符后就停止继续寻找
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[rune]bool) //创建字符集，也就是唯一的字符，方便后面判断是否存在
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r { // 如果左端和右端都没找到正常字符，那么表示该字符串没有正常字符
		return ""
	}
	return string(s[start : end+1])
}

func parseRESP(conn net.Conn) interface{} {
	r := bufio.NewReader(conn)
	line, err := r.ReadString('\n')
	if err != nil {
		return ""
	}

	cmdType := string(line[0])
	cmdTxt := strings.Trim(string(line[1:]), "\r\n")

	switch cmdType {
	case "*":
		count, _ := strconv.Atoi(cmdTxt)
		var data []string
		for i := 0; i < count; i++ {
			line, _ := r.ReadString('\n')
			cmd_txt := strings.Trim(string(line[1:]), "\r\n")
			c, _ := strconv.Atoi(cmd_txt)
			length := c + 2
			str := ""
			for length > 0 {
				block, _ := r.Peek(length)
				if length != len(block) {

				}
				r.Discard(length)
				str += string(block)
				length -= len(block)
			}

			data = append(data, strings.Trim(str, "\r\n"))
		}
		return data
	default:
		return cmdTxt
	}
}

func parseRedisConn(conn net.Conn) interface{} {
	r := bufio.NewReader(conn)
	line, err := r.ReadString('\n')
	if err != nil {
		return ""
	}

	cmdType := string(line[0])
	cmdTxt := strings.Trim(string(line[1:]), "\r\n")

	switch cmdType {
	case "$":
		count, _ := strconv.Atoi(cmdTxt)
		//var data []string
		buf := make([]byte, count)
		n, err := io.ReadFull(r, buf)
		if err != nil {
			fmt.Println(n, err.Error())
		} else {
			filename := "redis" + time.Now().Format("2006-01-02-15-04-05")
			f, err := os.Create(filename)
			defer f.Close()
			if err != nil {
				fmt.Println(n, err.Error())
			}
			_, err = f.Write(buf)
			if err != nil {
				fmt.Println(n, err.Error())
			}
			return buf
		}
	default:
		return cmdTxt
	}
	return cmdTxt
}
