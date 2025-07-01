import {createContext, type PropsWithChildren, useContext, useEffect, useState} from "react";
import type {User} from "@/types/user.ts";
import * as React from "react";
import {redirect, useNavigate} from 'react-router-dom';
import {jwtDecode} from "jwt-decode";

type AuthContext = {
    currentUser: User | null;
    handleLogout: () => Promise<void>;
}
const AuthContext = createContext<AuthContext | undefined>(undefined);

type AuthProviderProps = PropsWithChildren;

export const AuthProvider: React.FC<AuthProviderProps> = ({children}: AuthProviderProps) => {
    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        const getUser = async () => {
            const access_token = localStorage.getItem("access_token");
            if (access_token != null) {
                try {
                    const user = jwtDecode<User>(access_token);
                    if (user.full_name == "Hoang Bao") {
                        user.roles = "Client"
                    }
                    setCurrentUser(user);
                    redirect("/sign-in")
                    return;
                } catch (err) {
                    console.error("Invalid token", err);
                    localStorage.removeItem("access_token");
                }
            }
            navigate('/sign-in');
        }

        getUser();
    }, [navigate]);


    async function handleLogout() {
        localStorage.removeItem("access_token");
        setCurrentUser(null);
        navigate('/sign-in')
    }


    return <AuthContext.Provider
        value={{
            currentUser,
            handleLogout,
        }}
    >{children}</AuthContext.Provider>
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within a AuthProvider');
    }
    return context;
}