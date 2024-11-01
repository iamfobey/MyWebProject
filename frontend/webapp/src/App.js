import React, { useState, useEffect } from "react";
import { BrowserRouter, Route, Routes, useNavigate } from "react-router-dom";
import {BaseButton, AuthForm, AuthInput, AuthTitle, BaseDiv, AuthButton} from './App.styled';

const Home = () => {
    const navigate = useNavigate();
    return (
        <BaseDiv>
            <h1>Home</h1>

            <BaseButton onClick={() => {
                navigate("/account");
            }}>Account
            </BaseButton>
            <br/>
            <BaseButton onClick={() => {
                navigate("/auth");
            }}>Auth
            </BaseButton>
        </BaseDiv>
    );
}

const Auth = () => {
    const [authType, setAuthType] = useState("login");
    const [login, setLogin] = useState("");
    const [password, setPassword] = useState("");
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem("authToken");

        if (token) {
            navigate("/");
        }
    }, [navigate]);

    const handleSubmit = (e) => {
        e.preventDefault();

        const formData = new URLSearchParams();

        formData.append('authType', authType);
        formData.append('login', login);
        formData.append('password', password);

        fetch('http://localhost:8080/api/public/auth/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: formData.toString(),
        })
            .then(response => {
                if (!response.ok) throw new Error(response.statusText);
                return response.json();
            })
            .then(data => {
                localStorage.setItem("authToken", data.token);

                navigate("/");
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    };

    return (
        <BaseDiv>
            <AuthForm onSubmit={handleSubmit}>
                <AuthTitle>{authType === "login" ? "Login" : "Create your account"}</AuthTitle>
                <br/>
                <span>Username</span>
                <AuthInput type="text" name="login" value={login} onChange={(e) => setLogin(e.target.value)} required/>
                <br/>
                <span>Password</span>
                <AuthInput type="password" name="password" value={password}
                           onChange={(e) => setPassword(e.target.value)} required/>
                <br/>
                <AuthButton type="submit">{authType === "login" ? "Login" : "Create account"}</AuthButton>
                <p/>
                <hr/>
                <AuthButton type="button" onClick={() => setAuthType(authType === "login" ? "signup" : "login")}>
                    {authType === "login" ? "Create account" : "Login"}
                </AuthButton>
            </AuthForm>
        </BaseDiv>
    );
}

const Account = () => {
    const navigate = useNavigate();

    const handleLogout = () => {
        localStorage.removeItem("authToken");
        navigate("/");
    };

    return (
        <BaseDiv>
            <BaseButton onClick={handleLogout}>Logout</BaseButton>
        </BaseDiv>
    );
}

const App = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path={"/"} element={<Home />} />
                <Route path={"/auth"} element={<Auth />} />
                <Route path={"/account"} element={<Account />}></Route>
            </Routes>
        </BrowserRouter>
    );
};

export default App;