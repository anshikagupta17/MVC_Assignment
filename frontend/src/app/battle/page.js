"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/services/api";
import Link from "next/link";

export default function BattlePage() {
    const [opponent, setOpponent] = useState(null);
    const [loading, setLoading] = useState(false);
    const [battleResult, setBattleResult] = useState(null);
    const [troops, setTroops] = useState([]);
    const [selectedTroops, setSelectedTroops] = useState({});
    const [battlePhase, setBattlePhase] = useState("idle");
    const [battleMessage, setBattleMessage] = useState("");
    const [messageIndex, setMessageIndex] = useState(0);

    const battleMessages = [
        "Troops advancing",
        "Defenses fighting back",
        "Breaking through",
        "Townhall in DANGER",
        "Final assault underway",
    ];


    useEffect(() => {
        fetchTroops();
    }, []);

    async function fetchTroops() {
        try {
            const res = await apiFetch("/village/troops");
            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "Failed to load troops");
            }

            setTroops(data || []);
        } catch (err) {
            console.error(err);
        }
    }


    async function FindOpponent() {
        setLoading(true);
        setOpponent(null);
        setBattleResult(null);

        try {
            const res = await apiFetch("/battle/matchmake", {
                method: "POST",
            });

            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "No opponent found");
            }

            setOpponent(data);
        } catch (err) {
            console.error(err);
            alert(err.message);
        } finally {
            setLoading(false);
        }
    }

    function handleTroopQuantityChange(troopId, quantity) {
        setSelectedTroops((prev) => ({
            ...prev,
            [troopId]: quantity,
        }));
    }

    function buildTroopsToSend(selectedTroops) {
        const troopsToSend = [];

        for (const troopIdStr in selectedTroops) {
            const quantity = selectedTroops[troopIdStr];

            if (quantity > 0) {
                troopsToSend.push({
                    troop_id: Number(troopIdStr),
                    quantity: quantity,
                });
            }
        }

        return troopsToSend;
    }

    function advanceBattleMessage() {
        setMessageIndex(function nextIndex(prevIndex) {
            const newIndex = (prevIndex + 1) % battleMessages.length;
            setBattleMessage(battleMessages[newIndex]);
            return newIndex;
        });
    }

    function waitMilliseconds(ms) {
        return new Promise(function startTimer(resolve) {
            setTimeout(resolve, ms);
        });
    }

    async function attackOpponent() {
        if (!opponent) return;

        const troopsToSend = buildTroopsToSend(selectedTroops);

        if (troopsToSend.length === 0) {
            alert("Select at least one troop to attack with");
            return;
        }

        setBattlePhase("fighting");
        setMessageIndex(0);
        setBattleMessage(battleMessages[0]);

        const messageInterval = setInterval(advanceBattleMessage, 6000);

        try {
            const res = await apiFetch("/battle/attack", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    defender_id: opponent.village_id,
                    troops: troopsToSend,
                }),
            });

            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "Attack failed");
            }

            clearInterval(messageInterval);
            setBattleResult(data);
            setBattlePhase("result");
            setOpponent(null);
            setSelectedTroops({});
            await fetchTroops();
        } catch (err) {
            clearInterval(messageInterval);
            setBattlePhase("idle");
            console.error(err);
            alert(err.message);
        }
    }

    return (
        <div className="p-6 space-y-6"
            style={{backgroundImage: `url('/assets/logo.png')`,
                backgroundRepeat: "no-repeat",

            }}
        >
            <div className="flex justify-between items-center"
        >
                <h1 className="text-3xl font-bold">Battle</h1>
                <Link href="/village">
                    <button className="border px-4 py-2">
                        Back to Village
                    </button>
                </Link>
            </div>
            <button
                onClick={FindOpponent}
                disabled={loading}
                className="border px-4 py-2 bg-red-500 text-white"
            >
                {loading ? "Searching..." : "Find Opponent"}
            </button>

            {opponent && battlePhase === "idle" && (
                <div className="border p-4 space-y-4">
                    <div>
                        <h2 className="text-xl font-semibold">Opponent Found</h2>
                        <p>Village ID: {opponent.village_id}</p>
                        <p>Townhall Level: {opponent.townhall_level}</p>
                        <p>Trophies: {opponent.trophies}</p>
                        <p>Gold: {opponent.gold}</p>
                        <p>Elixir: {opponent.elixir}</p>
                    </div>

                    <div>
                        <h3 className="text-lg font-semibold mb-2">Select Troops to Attack</h3>

                        {troops.length === 0 ? (
                            <p>You have no troops. Train some first.</p>
                        ) : (
                            troops.map((t) => (
                                <div key={t.troop_id} className="flex items-center justify-between border p-2 mb-2">
                                    <div>
                                        <p className="font-semibold">{t.name}</p>
                                        <p className="text-sm">Owned: {t.quantity}</p>
                                    </div>
                                    <input
                                        type="number"
                                        min="0"
                                        max={t.quantity}
                                        value={selectedTroops[t.troop_id] || 0}
                                        onChange={(e) =>
                                            handleTroopQuantityChange(t.troop_id, Number(e.target.value))
                                        }
                                        className="border px-2 py-1 w-20"
                                    />
                                </div>
                            ))
                        )}
                    </div>

                    <button
                        onClick={attackOpponent}
                        className="border px-4 py-2 bg-green-500 text-white"
                    >
                        Attack
                    </button>
                </div>
            )}

            {battlePhase === "fighting" && (
                <div className="fixed inset-0 bg-black flex flex-col items-center justify-center z-50 text-white">
                    <div className="relative w-[600px] h-[600px] border-4 border-red-700 bg-[#3F704D] overflow-hidden">
                        {opponent?.buildings?.map((b) => (
                            <div
                                key={b.id}
                                className="absolute bg-black border"
                                style={{
                                    left: `${b.x * 10}px`,
                                    top: `${b.y * 10}px`,
                                    width: `${b.size_x * 10}px`,
                                    height: `${b.size_y * 10}px`,
                                }}
                            />
                        ))}
                    </div>

                    <p className="mt-6 text-xl animate-pulse">{battleMessage}</p>
                </div>
            )}

            {battlePhase === "result" && battleResult && (
                <div className="border p-4 space-y-2">
                    <h2 className="text-xl font-semibold">Battle Result</h2>
                    <p>Stars: {battleResult.stars}</p>
                    <p>Destruction: {battleResult.destruction}%</p>
                    <p>Loot Gold: {battleResult.loot_gold}</p>
                    <p>Loot Elixir: {battleResult.loot_elixir}</p>
                    <p>Trophy Change: {battleResult.trophy_change}</p>
                </div>
            )}
        </div>
    );
}