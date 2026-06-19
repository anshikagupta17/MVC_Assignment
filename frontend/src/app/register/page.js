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
        <div className="p-8">
            <h1 className="text-3xl font-bold mb-6">Register</h1>

            <form onSubmit={handleSubmit} className="flex flex-col gap-4 max-w-sm">
                <input
                    type="text"
                    placeholder="Username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    className="border p-2"
                />

                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="border p-2"
                />

                <button type="submit" className="border p-2">
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
    );
}