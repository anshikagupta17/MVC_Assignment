"use client";
import { BUILDING_NAMES } from "@/constants/buildings";

export default function Building({
    building,
    onSelect,
    isSelected,
}) {
    function handleClick(e) {
        e.stopPropagation();
        onSelect(building);
    }
    return (
        <div
            onClick={handleClick}
            className={`absolute w-16 h-16 border flex flex-col items-center justify-center text-xs text-center shadow cursor-pointer hover:scale-105 transition
                ${building.upgrade_ends_at ? "opacity-60 border-red-400" : "bg-black"}
                ${isSelected ? "ring-2 ring-yellow-300" : ""}
            `}
            style={{
                left: `${building.x * 20}px`,
                top: `${building.y * 20}px`,
                width: `${building.size_x*20}px`,
                height: `${building.size_y*20}px`,
            }}
        >
            {building.upgrade_ends_at && (
                <div className="absolute -top-4 text-[10px] text-red-400 bg-black px-1 rounded">
                    Upgrading
                </div>
            )}

            <div>{BUILDING_NAMES[building.building_id]}</div>

            <div>Lv {building.level}</div>
        </div>
    );
}