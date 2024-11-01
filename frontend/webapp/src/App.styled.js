import styled from '@emotion/styled';

export const BaseDiv = styled.div`
    font-family: Arial, sans-serif;
`;

export const AuthForm = styled.form`
    max-width: 400px;
    margin: 0 auto;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 8px;
    background-color: #f9f9f9;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
`;

export const AuthInput = styled.input`
    width: 100%;
    padding: 10px;
    margin: 10px -10px;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 16px;

    &:focus {
        border-color: #007bff;
        outline: none;
    }
`;

export const AuthTitle = styled.p`
    text-align: center;
    font-size: 20px;
`;

export const BaseButton = styled.button`
    width: auto;
    padding: 10px;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    cursor: pointer;

    &:hover {
        background-color: #0056b3;
    }
`;

export const AuthButton = styled(BaseButton)`
    width: 100%;
`;