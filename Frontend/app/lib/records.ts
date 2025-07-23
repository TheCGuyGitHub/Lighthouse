
export interface Server {
    id: string;
    name: string;
    hostname: string;
    port: string;
    created: string;
    updated: string;
    state: string;
    quickdata?: {
        ram?: {
            used?: number;
            total?: number;
        };
        cpu?: {
            usage?: number;
        };
        disk?: {
            used?: number;
            total?: number;
        };
        temp?: number;
        agent_version?: string;
    };
}