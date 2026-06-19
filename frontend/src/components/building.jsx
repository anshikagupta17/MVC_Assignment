"use client";
import { BUILDING_NAMES } from "@/constants/buildings";
import Image from "next/image";

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
            className={`absolute cursor-pointer hover:scale-105 transition
                ${isSelected ? "ring-2 ring-yellow-300 rounded" : ""}
            `}
            style={{
                left: `${building.x * 20}px`,
                top: `${building.y * 20}px`,
                width: `${building.size_x * 20*1.2}px`,
                height: `${building.size_y * 20*1.2}px`,
            }}
        >
            <Image
                src={`/assets/buildings/${building.building_id}.png`}
                alt={BUILDING_NAMES[building.building_id]}
                fill
                sizes="100px"
                className={building.upgrade_ends_at ? "opacity-60" : ""}
            ></Image>
        </div>
    );
}