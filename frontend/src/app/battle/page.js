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
    const [errorMsg, setErrorMsg] = useState("");

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

    function showError(rawMessage, fallbackMessage) {
        const displayMessage = fallbackMessage || rawMessage;
        setErrorMsg(displayMessage);
        setTimeout(() => setErrorMsg(""), 5000);
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
            showError(err.message, "Couldn't find an opponent right now. Try again in a bit.");
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
            showError(null, "Pick at least one troop before attacking.");
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

            await waitMilliseconds(30000);

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
            showError(err.message, "The attack couldn't go through. Try again.");
        }
    }

    function getSelectedTroopList() {
        const list = [];

        for (const troopId in selectedTroops) {
            const quantity = selectedTroops[troopId];
            if (quantity > 0) {
                const troop_Id = Number(troopId);
                const troopData = troops.find(function findTroop(t) {
                    return t.troop_id === troopId;
                });

                if (troopData) {
                    list.push({
                        troop_id: troop_Id,
                        name: troopData.name,
                        quantity: quantity,
                    });
                }
            }
        }

        return list;
    }

    return (
        <div className="min-h-screen bg-[#1a0a00] text-white p-6">
            {errorMsg && (
                <div className="fixed top-6 left-1/2 -translate-x-1/2 bg-red-600 text-white px-4 py-2 rounded shadow-lg z-50">
                    {errorMsg}
                </div>
            )}
            <div className="flex justify-between items-center mb-6">
                <Link href="/village">
                    <button className="bg-[#3d1f00] border border-[#8B5E3C] text-[#f5c96e] px-4 py-2 rounded hover:bg-[#5a2e00]">
                        Back to Village
                    </button>
                </Link>
            </div>

            {battlePhase === "idle" && (
                <div className="flex flex-col items-center gap-6">
                    <button
                        onClick={FindOpponent}
                        disabled={loading}
                        className="bg-red-700 hover:bg-red-600 text-white text-xl font-bold px-10 py-4 rounded-xl border-2 border-red-400 shadow-lg"
                    >
                        {loading ? "Searching..." : "Find Opponent"}
                    </button>
                    {opponent && (
                    <div className="w-full max-w-5xl flex flex-col lg:flex-row gap-6">

                            <div className="flex-1 bg-[#2a1200] border-2 border-[#8B5E3C] rounded-xl p-4">
                                <h2 className="text-xl font-bold text-[#f5c96e] mb-3 text-center">
                                    Enemy Village
                                </h2>

                                <div className="flex gap-6 justify-center mb-4 text-sm">
                                    <span>Townhall {opponent.townhall_level}</span>
                                    <span>Trophies {opponent.trophies}</span>
                                    <span>Gold {opponent.gold}</span>
                                    <span>ELixir {opponent.elixir}</span>
                                </div>

                                <div
                                    className="relative mx-auto overflow-hidden border border-[#8B5E3C]"
                                    style={{
                                        width: "500px",
                                        height: "500px",
                                        backgroundColor: "#3F704D",
                                        backgroundImage: `
                                            linear-gradient(rgba(255,255,255,0.08) 1px, transparent 1px),
                                            linear-gradient(90deg, rgba(255,255,255,0.08) 1px, transparent 1px)
                                        `,
                                        backgroundSize: "10px 10px",
                                    }}
                                >
                                    {opponent.buildings?.map((b) => (
                                        <div
                                            key={b.id}
                                            className="absolute"
                                            style={{
                                                left: `${b.x * 10}px`,
                                                top: `${b.y * 10}px`,
                                                width: `${b.size_x * 10}px`,
                                                height: `${b.size_y * 10}px`,
                                            }}
                                        >
                                            <img
                                                src={`/assets/buildings/${b.building_id}.png`}
                                                alt=""
                                                style={{
                                                    width: "100%",
                                                    height: "100%",
                                                    objectFit: "contain",
                                                }}
                                            />
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="flex-1 bg-[#2a1200] border-2 border-[#8B5E3C] rounded-xl p-4">
                                <h2 className="text-xl font-bold text-[#f5c96e] mb-3 text-center">
                                    Select Troops
                                </h2>

                                <div className="grid grid-cols-2 gap-3 mb-4">
                                    {troops.filter(function hasTroops(t) {
                                        return t.quantity > 0;
                                    }).map((t) => (
                                        <div
                                            key={t.troop_id}
                                            className={`relative bg-[#3d1f00] border-2 rounded-xl p-3 flex flex-col items-center gap-1 cursor-pointer transition-all ${
                                                selectedTroops[t.troop_id] > 0
                                                    ? "border-yellow-400 shadow-yellow-400 shadow-md"
                                                    : "border-[#8B5E3C]"
                                            }`}
                                        >
                                            <p className="text-[#f5c96e] font-bold text-sm tracking-wide uppercase">
                                                {t.name}
                                            </p>

                                            <img
                                                src={`/assets/troops/${t.troop_id}.png`}
                                                alt={t.name}
                                                className="w-16 h-16 object-contain"
                                            />

                                            <div className="flex items-center gap-1 text-xs text-yellow-300">
                                                <span>Owned: {t.quantity}</span>
                                            </div>

                                            <input
                                                type="number"
                                                min="0"
                                                max={t.quantity}
                                                value={selectedTroops[t.troop_id] || 0}
                                                onChange={(e) =>
                                                    handleTroopQuantityChange(t.troop_id, Number(e.target.value))
                                                }
                                                className="w-16 text-center bg-[#1a0a00] border border-[#8B5E3C] text-white rounded px-1 py-0.5 text-sm"
                                            />
                                        </div>
                                    ))}
                                </div>

                                <button
                                    onClick={attackOpponent}
                                    className="w-full bg-green-700 hover:bg-green-600 text-white font-bold py-3 rounded-xl border-2 border-green-400 text-lg"
                                >
                                    Attack
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            )}

            {battlePhase === "fighting" && (
                <div className="fixed inset-0 bg-black flex flex-col items-center justify-center z-50">
                    <p className="text-4xl font-bold text-[#f5c96e] animate-pulse tracking-widest" style={{ fontFamily: "serif" }}>
                        {battleMessage}
                    </p>
                </div>
            )}

            {battlePhase === "result" && battleResult && (
                <div className="max-w-md mx-auto bg-[#2a1200] border-2 border-[#8B5E3C] rounded-xl p-6 text-center">
                    <h2 className="text-2xl font-bold text-[#f5c96e] mb-4" style={{ fontFamily: "serif" }}>
                        Battle Result
                    </h2>

                    <div className="flex justify-center gap-2 mb-4">
                        {[1, 2, 3].map((star) => (
                            <span
                                key={star}
                                className={`text-4xl ${star <= battleResult.stars ? "text-yellow-400" : "text-gray-600"}`}
                            >
                                ★
                            </span>
                        ))}
                    </div>

                    <div className="space-y-2 text-lg">
                        <p>Destruction: <span className="text-[#f5c96e] font-bold">{battleResult.destruction_percent}%</span></p>
                        <p>Gold looted: <span className="text-yellow-300 font-bold">{battleResult.loot_gold}</span></p>
                        <p>Elixir looted: <span className="text-purple-300 font-bold">{battleResult.loot_elixir}</span></p>
                        <p>Trophy change: <span className={`font-bold ${battleResult.trophy_change >= 0 ? "text-green-400" : "text-red-400"}`}>
                            {battleResult.trophy_change >= 0 ? "+" : ""}{battleResult.trophy_change}
                        </span></p>
                    </div>

                    <button
                        onClick={FindOpponent}
                        className="mt-6 bg-red-700 hover:bg-red-600 text-white font-bold px-8 py-3 rounded-xl border-2 border-red-400"
                    >
                        Battle Again
                    </button>
                </div>
            )}
        </div>
    );
}