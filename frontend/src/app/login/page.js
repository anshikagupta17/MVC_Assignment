"use client";

import { useState } from "react";
import { apiFetch } from "@/services/api";
import { useAuth } from "@/context/auth_context";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const { login } = useAuth();
   const router = useRouter();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();

    try {
        const response = await apiFetch("/login", {
            method: "POST",
            body: JSON.stringify({ username, password }),
        });

        const text = await response.text();

        if (!response.ok) {
            alert(text || "Login failed");
            return;
        }

        const data = JSON.parse(text);
        login(data.token);
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
          Login
        </h1>

        <form
          onSubmit={handleSubmit}
          className="flex flex-col gap-4"
        >
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) =>
              setUsername(e.target.value)
            }
            className="border p-2 rounded"
          />

          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) =>
              setPassword(e.target.value)
            }
            className="border p-2 rounded"
          />

          <button
            type="submit"
            className="border p-2 rounded bg-black text-white"
          >
            Login
          </button>
        </form>
      </div>
    </div>
  );
}