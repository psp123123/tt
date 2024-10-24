worker_processes auto;
error_log /var/log/nginx/error.log warn;  # 设置更低的日志级别，避免过多日志写入
pid /run/nginx.pid;

include /usr/share/nginx/modules/*.conf;

events {
    multi_accept on;  # 让每个 worker 进程尽可能多地接受新连接
    use epoll;  # 在 Linux 上使用高效的事件模型
}

http {
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';
    access_log /var/log/nginx/access.log main;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    keepalive_requests 100;  # 限制每个连接的请求数量，防止长时间占用资源

    types_hash_max_size 2048;
    client_max_body_size 16m;  # 限制上传文件的大小，增强安全性

    include /etc/nginx/mime.types;  # 确保包含 mime.types 文件以处理常见的文件类型

    gzip on;  # 开启 gzip 压缩，提高传输效率
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
    gzip_min_length 1024;  # 只压缩大于 1KB 的文件
    gzip_comp_level 5;  # 压缩级别（1-9），推荐 5 以获得平衡的性能和压缩率

    include /etc/nginx/conf.d/*.conf;

    server {
        listen 80 default_server;
        listen [::]:80 default_server;
        server_name _;

        root /data/prod/LayuiMini;
        index index.html;

        include /etc/nginx/default.d/*.conf;

        location / {
            try_files $uri $uri/ /index.html;  # 优化静态文件查找逻辑，提高响应效率
            root /data/prod/LayuiMini;
            index index.html;
        }

        error_page 404 /404.html;
        location = /404.html {
            internal;  # 将 404 页面设为 internal，防止外部访问
        }

        error_page 500 502 503 504 /50x.html;
        location = /50x.html {
            internal;  # 同样将 50x 页面设为 internal
        }

        # 安全增强配置
        add_header X-Content-Type-Options nosniff;
        add_header X-Frame-Options SAMEORIGIN;
        add_header X-XSS-Protection "1; mode=block";
    }
}