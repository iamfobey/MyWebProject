import React, { useState, useEffect } from "react";
import { BrowserRouter, Route, Routes, useNavigate } from "react-router-dom";
import {BaseButton, AuthForm, AuthInput, AuthTitle, BaseDiv, AuthButton} from './App.styled';

const Home = () => {
    const navigate = useNavigate();
    return (
        <BaseDiv>
            <h1>Home</h1>

            {localStorage.getItem("authToken") !== null &&
                <BaseButton onClick={() => {
                    navigate("/account");
                }}>Account
                </BaseButton>
            }
            {localStorage.getItem("authToken") === null &&
                <BaseButton onClick={() => {
                    navigate("/auth");
                }}>Auth
                </BaseButton>
            }
        </BaseDiv>
    );
}

const Auth = () => {
    const [authType, setAuthType] = useState("login");
    const [login, setLogin] = useState("");
    const [password, setPassword] = useState("");
    const [loginError, setLoginError] = useState("");
    const [passwordError, setPasswordError] = useState("");

    const navigate = useNavigate();

    const validateLogin = (value) => {
        const loginRegex = /^[a-zA-Z0-9_]{3,20}$/;
        if (!loginRegex.test(value)) {
            setLoginError("Логин должен содержать только буквы, цифры, подчеркивания и быть длиной от 3 до 20 символов.");
        } else {
            setLoginError("");
        }
    };

    const validatePassword = (value) => {
        const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,48}$/;
        if (!passwordRegex.test(value)) {
            setPasswordError("Пароль должен содержать от 8 до 48 символов, включая заглавные и строчные буквы, цифру и специальный символ.");
        } else {
            setPasswordError("");
        }
    };

    const handleLoginChange = (e) => {
        const value = e.target.value;
        setLogin(value);
        validateLogin(value);
    };

    const handlePasswordChange = (e) => {
        const value = e.target.value;
        setPassword(value);
        validatePassword(value);
    };

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

        fetch('http://192.168.0.106:8080/api/public/auth/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: formData.toString(),
        })
            .then(response => {
                if (!response.ok) throw new Error(response.statusMessage + " " + response.statusText);
                return response.json();
            })
            .then(data => {
                localStorage.setItem("authToken", data.token);
                setLogin('');
                setPassword('');
                navigate("/");
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    };

    const handleChangeType = () => {
        setAuthType(authType === "login" ? "signup" : "login");
        setLogin('');
        setPassword('');
        setLoginError('');
        setPasswordError('');
    }

    return (
        <BaseDiv>
            <AuthForm onSubmit={handleSubmit}>
                <AuthTitle>{authType === "login" ? "Login" : "Create your account"}</AuthTitle>
                <br/>
                <span>Username</span>
                <AuthInput
                    type="text"
                    name="login"
                    value={login}
                    onChange={handleLoginChange}
                    required
                />
                {authType === "signup" && loginError && <p style={{ color: 'red' }}>{loginError}</p>}
                <br/>
                <span>Password</span>
                <AuthInput
                    type="password"
                    name="password"
                    value={password}
                    onChange={handlePasswordChange}
                    required
                />
                {authType === "signup" && passwordError && <p style={{ color: 'red' }}>{passwordError}</p>}
                <br/>
                <AuthButton type="submit">
                    {authType === "login" ? "Login" : "Create account"}
                </AuthButton>
                <p/>
                <hr/>
                <AuthButton type="button" onClick={handleChangeType}>
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