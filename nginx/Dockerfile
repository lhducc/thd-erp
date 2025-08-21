# FROM nginx:1.27-alpine
FROM 192.168.1.76:5001/product/erp/erp/nginx:1.27

COPY --from=builder /src/dist /usr/share/nginx/html

# Create SSL directory and copy certificates
RUN mkdir -p /etc/nginx/ssl
COPY ssl/thd.erp.com.crt /etc/nginx/ssl/server.crt
COPY ssl/thd.erp.com.key /etc/nginx/ssl/server.key

# Set proper permissions for SSL files and ensure nginx user ownership
RUN chmod 644 /etc/nginx/ssl/server.crt && \
    chmod 600 /etc/nginx/ssl/server.key && \
    chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /etc/nginx/ssl && \
    touch /var/run/nginx.pid && \
    chown nginx:nginx /var/run/nginx.pid

# Override default nginx.conf
COPY nginx.conf /etc/nginx/nginx.conf

USER nginx

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]