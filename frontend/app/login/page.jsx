'use client'
import Image from "next/image";
import { useState } from "react";
import Button from "../components/Button";
import { useRouter } from "next/navigation";
import Cookies from "js-cookie";
import getLocation from "../helpers/getLocation";

export default function Login() {

    const [name, setName] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const router = useRouter()

    async function handleSubmit(e) {
        e.preventDefault();

        const formData = new FormData();
        formData.append("name", name);
        formData.append("password", password);

        const req = await fetch(`${getLocation()}/api/login`, {
            method: 'post',
            body: formData,
        });
        
        if (req.ok) {
            router.push('/pipeline');
        } else {
            setError('Invalid password/username');
        }

    }

    return (
        <div className="w-full min-h-screen flex justify-center items-center">
            <div className="border-2 border-gray-85 rounded-xl p-3">
                <form className="flex flex-col gap-2 mb-3"
                    onSubmit={handleSubmit}
                >
                    <h1 className="text-center font-bold text-2xl">
                        Sign in
                    </h1>
                    <label className="flex flex-col">
                        <span className="mb-1">
                            Name
                        </span>
                        <input
                            className="border rounded p-1.5"
                            id="name"
                            value={name}
                            required
                            onChange={(e) => {setName(e.currentTarget.value), setError('')}}
                        />
                    </label>
                    <label className="flex flex-col">
                        <span className="mb-1">
                            Password
                        </span>
                        <input
                        className="border rounded p-1.5"
                            id="password"
                            type="password"
                            value={password}
                            required
                            onChange={(e) => {setPassword(e.currentTarget.value), setError('')}}
                        />
                    </label>
                    {!error ? '' : (
                        <span className="my-1 text-red-400 text-center">
                            {error}
                        </span>
                    )}
                    <Button
                        className="py-2"
                        type="Submit"
                    >
                        Login
                    </Button>
                </form>
            </div>
        </div>

        
    );
}
