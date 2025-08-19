# Cấu hình MinIO HTTPS cho ERP System

## Tổng quan
Tài liệu này mô tả các thay đổi được thực hiện để kích hoạt HTTPS cho MinIO và đảm bảo các URL object có thể truy cập từ trình duyệt bên ngoài.

## Thay đổi chính

### 1. Docker Compose Configuration (`docker-compose.uat.yml`)
- **MinIO Service**: Đã cấu hình SSL certificates cho MinIO
- **SSL Certificates**: Mount certificates vào `/root/.minio/certs/`
  ```yaml
  volumes:
    - ../ssl/thd.erp.com.crt:/root/.minio/certs/public.crt:ro
    - ../ssl/thd.erp.com.key:/root/.minio/certs/private.key:ro
  ```
- **Environment Variables**: 
  - `MINIO_SERVER_URL`: Sử dụng HTTPS endpoint
  - `MINIO_HEALTH_URL`: Sử dụng HTTPS cho health check

### 2. Backend Configuration

#### MinIO Config (`backend/config/minIO_config.go`)
- **Dual Client Setup**: Tạo 2 MinIO clients
  - `MinioClient`: Cho operations nội bộ (sử dụng Docker service name `minio:9000`)
  - `MinioURLClient`: Cho tạo presigned URLs (sử dụng `localhost:9000` với custom transport)

- **Custom Transport**: Redirect `localhost:9000` requests đến `172.18.0.1:9000` (Docker host IP)
  ```go
  DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
      if strings.Contains(addr, "localhost:9000") {
          addr = "172.18.0.1:9000"
      }
      return customDialer.DialContext(ctx, network, addr)
  }
  ```

- **SSL Configuration**: Skip SSL verification cho internal communication
  ```go
  TLSClientConfig: &tls.Config{InsecureSkipVerify: true}
  ```

#### Bucket Operations (`backend/pkg/minIO/bucket.go`)
- **Presigned URL Generation**: Sử dụng `MinioURLClient` để tạo URLs với hostname đúng
- **Public Bucket Policy**: Tự động set bucket policy cho phép public read access
- **Error Handling**: Improved error handling cho SSL và signature issues

## Cách hoạt động

### 1. Upload Flow
1. Backend sử dụng `MinioClient` (internal) để upload files qua `minio:9000`
2. SSL verification được skip cho internal communication
3. Bucket policy được tự động set để cho phép public access

### 2. URL Generation Flow
1. Backend sử dụng `MinioURLClient` để generate presigned URLs
2. Client connect tới `localhost:9000` nhưng được redirect tới `172.18.0.1:9000`
3. Signature được tạo với hostname `localhost:9000`
4. URL trả về có format: `https://localhost:9000/bucket/object?signature=...`

### 3. Browser Access
1. Browser truy cập URL với `https://localhost:9000`
2. MinIO server response với SSL certificate
3. Object được serve thành công

## Lưu ý quan trọng cho Developers

### 1. SSL Certificates
- Đảm bảo files `thd.erp.com.crt` và `thd.erp.com.key` tồn tại trong thư mục `ssl/`
- Certificates phải valid và trust được bởi browser
- Restart containers sau khi update certificates

### 2. Network Configuration
- Backend container cần access được host network qua `172.18.0.1`
- Không thay đổi Docker network mode về `host` (sẽ break database connectivity)
- Port 9000 và 9001 phải được expose từ MinIO container

### 3. Environment Variables
- `MINIO_ENDPOINT`: Luôn sử dụng Docker service name cho internal communication
- `MINIO_PUBLIC_ENDPOINT`: Luôn sử dụng full HTTPS URL cho external access
- Không sử dụng IP addresses trong config production

### 4. Troubleshooting

#### Error: "SignatureDoesNotMatch"
- Kiểm tra hostname trong presigned URL có match với client endpoint
- Verify access keys và secret keys
- Đảm bảo clock sync giữa containers

#### Error: "NoSuchKey"
- File có thể đã bị xóa bởi lifecycle policy
- Verify object name và bucket name
- Check bucket permissions

#### Error: "Access Denied"
- Bucket policy chưa được set correctly
- Verify SSL certificates
- Check MinIO admin permissions

### 5. Testing Commands
```bash
# Test MinIO accessibility
curl -k https://localhost:9000/

# Check container network
docker exec erp-backend ip route show default

# Restart services
docker compose -f deploys/docker-compose.uat.yml restart backend minio

# View logs
docker compose -f deploys/docker-compose.uat.yml logs -f backend
docker compose -f deploys/docker-compose.uat.yml logs -f minio
```

## Security Considerations

1. **SSL Certificate Management**: 
   - Sử dụng proper CA-signed certificates cho production
   - Rotate certificates thường xuyên
   - Monitor certificate expiry

2. **Access Control**:
   - Bucket policies chỉ cho phép read access
   - Admin credentials được bảo vệ trong environment variables
   - Network isolation giữa services

3. **Data Protection**:
   - All communication encrypted với HTTPS/TLS
   - Presigned URLs có expiry time (default 15 minutes)
   - File lifecycle policies để auto-cleanup

## Tác giả
- Được cấu hình bởi Devops Team
- Ngày: 18/08/2025
- Version: 1.0
