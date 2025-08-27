import axios from "axios";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    // headers: {
    //     "Content-Type": "application/json",
    // },
});

// Request interceptor
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem("access_token");
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        if (config.data instanceof FormData) {
            console.log("form")
            delete config.headers["Content-Type"];
        } else {
            config.headers["Content-Type"] = "application/json";
        }
        // Log request details
        console.log('Request:', {
            url: `${config.baseURL}${config.url}`,
            method: config.method?.toUpperCase(),
            headers: config.headers,
            params: config.params,
            data: config.data
        });

        return config;
    },
    (error) => {
        // Log request error
        console.error('Request Error:', error);
        return Promise.reject(error);
    }
);

// Response interceptor
api.interceptors.response.use(
    (response) => {
        // Log successful response
        console.log('Response:', {
            status: response.status,
            statusText: response.statusText,
            data: response.data,
            headers: response.headers,
            config: {
                url: response.config.url,
                method: response.config.method?.toUpperCase()
            }
        });

        return response;
    },
    (error) => {
        if (error.response) {
            // Log error response
            console.error('Response Error:', {
                status: error.response.status,
                statusText: error.response.statusText,
                data: error.response.data,
                headers: error.response.headers,
                config: {
                    url: error.config.url,
                    method: error.config.method?.toUpperCase()
                }
            });

            if (error.response.status === 401) {
                localStorage.removeItem("access_token");
                window.location.href = "/login";
            }
        } else {
            // Log error without response (network error, etc.)
            console.error('Error:', error.message);
        }

        return Promise.reject(error);
    }
);

export default api;