import type { ExtendedClient } from "@/query/client";
import type { HeaderAction, HeaderInput } from "./types";

function triggerDownload(filename: string, content: string) {
  const blob = new Blob([content], { type: "text/plain" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}

function generateAndDownloadAnsibleInventory({ rows: clients }: { rows: ExtendedClient[] }): void {
  const lines: string[] = ["[contest]"];
  clients.forEach((client) => lines.push(client.ip))
  triggerDownload("inventory.ini", lines.join("\n"));
}

function generateAndDownloadClusterSSHInventory({ rows: clients }: { rows: ExtendedClient[] }): void {
  const ips = clients.filter(item => item.ip).map(item => item.ip).join(" ");
  const content = `contest ${ips}`;
  triggerDownload("clusters", content);
}

export const generateAnsibleAction: HeaderAction<ExtendedClient> = {
  label: "Ansible",
  onClick: generateAndDownloadAnsibleInventory,
  disabled: ({ rows }: HeaderInput<ExtendedClient>) => rows.length === 0
}

export const generateClusterSSHAction: HeaderAction<ExtendedClient> = {
  label: "Cluster SSH",
  onClick: generateAndDownloadClusterSSHInventory,
  disabled: ({ rows }: HeaderInput<ExtendedClient>) => rows.length === 0
}
