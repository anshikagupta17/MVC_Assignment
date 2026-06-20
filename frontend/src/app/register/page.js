"use client";

import { useState } from "react";
import { apiFetch } from "@/services/api";
import { useAuth } from "@/context/auth_context";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function RegisterPage() {
    const { login } = useAuth();
    const router = useRouter();

    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    async function handleSubmit(e) {
        e.preventDefault();

        try {
            const registerRes = await apiFetch("/register", {
                method: "POST",
                body: JSON.stringify({ username, password }),
            });

            const registerData = await registerRes.json();

            if (!registerRes.ok) {
                alert(registerData.message || "Registration failed");
                return;
            }

            const loginRes = await apiFetch("/login", {
                method: "POST",
                body: JSON.stringify({ username, password }),
            });

            const loginData = await loginRes.json();

            if (!loginRes.ok) {
                alert(loginData.message || "Auto-login failed, please log in manually");
                router.push("/login");
                return;
            }

            login(loginData.token);
            router.push("/village");
        } catch (err) {
            console.error(err);
            alert("Server error");
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center"
            style={{backgroundImage:`url('/assets/battle.png')`}}
        >
            <div className="bg-[white]/60 text-black p-8 rounded-xl shadow-lg w-full max-w-sm">
                <h1 className="text-3xl font-bold mb-6 text-center">
                    Register
                </h1>

            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                <input
                    type="text"
                    placeholder="Username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    className="border p-2 rounded"
                />

                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="border p-2 rounded"
                />

                <button type="submit" className="border p-2 rounded bg-black text-white">
                    Register
                </button>
            </form>

            <p className="mt-4">
                Already have an account?{" "}
                <Link href="/login" className="underline">
                    Login
                </Link>
            </p>
            </div>
        </div>
    );
}