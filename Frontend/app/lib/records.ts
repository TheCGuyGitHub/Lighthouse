import type { ZodFloat64 } from "zod";

export interface Server {
    id: string;
    name: string;
    hostname: string;
    port: string;
    created: string;
    updated: string;
    state: string;
    quickdata?: {
        cpu_usage?: number;
        memory_cached?: number;
        memory_used?: number;
        memory_total?: number;
        disk_used?: number;
        disk_total?: number;
        disk_io?: {
            read: number;
            write: number;
        };
        bandwidth?: {
            download: ZodFloat64;
            upload: ZodFloat64;
        };
        temp_sensors?: Array<{
            name: string;
            value: number;
        }>;
        agent_version?: string;
    };
}

// Data from cpu usage, memory cached, memory used, memory total, disk used, disk total, disk io, bandwith, all temp sensors, agent version, gpu, docker, proxmox.
// NOTE: Stuff like temp, gpu, docker, proxmox are not always available, so they should be optional.
export interface ServerData {
    server_id: string;
    created: string;
    stats: {
        cpu_usage?: number;
        memory_cached?: number;
        memory_used?: number;
        memory_total?: number;
        disk_used?: number;
        disk_total?: number;
        disk_io?: {
            read: number;
            write: number;
        };
        bandwidth?: {
            download: ZodFloat64;
            upload: ZodFloat64;
        };
        temp_sensors?: Array<{
            name: string;
            value: number;
        }>;
        gpu?: {
            model: string;
            usage: number;
            memory_used: number;
            memory_total: number;
        };
        docker?: {
            containers_running: number;
            containers_total: number;
        };
        proxmox?: {
            nodes: Array<{
                name: string;
                status: string;
                cpu_usage: number;
                memory_used: number;
                memory_total: number;
            }>;
        };
    }
    type?: "1m" | "5m" | "15m" | "30m" | "60m" | "120m" | "240m";
}