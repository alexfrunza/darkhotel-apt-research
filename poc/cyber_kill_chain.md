# Linux Cyber Kill Chain & MITRE ATT&CK Breakdown

Darkhotel P2P intrusion model in a **Linux environment**. Instead of a Windows executable (`AxDecrypt.exe` / `Catch.exe`), the lure relies on a trojanized Linux binary bundled with encrypted archives, paired with a local privilege escalation (LPE) exploit targeting `sudo` to establish a root-level `supervisord` daemon.

Below is the threat report detailing the Cyber Kill Chain stages and corresponding **MITRE ATT&CK for Linux** Tactics, Techniques, and Procedures (TTPs).

---

## 1. Reconnaissance
The adversary profiles Linux user communities (such as developers, researchers, or open-source enthusiasts) on P2P networks and torrent trackers to select a high-demand lure, such as a large software dataset, media archive, or leaked documentation bundle.
* **MITRE Tactic:** Reconnaissance ([`TA0043`](https://attack.mitre.org/tactics/TA0043/))
* **MITRE Technique:** **Gather Victim Identity Information** ([`T1589`](https://attack.mitre.org/techniques/T1589/)) / **Search Open Technical Databases** ([`T1596`](https://attack.mitre.org/techniques/T1596/))
* **Explanation:** Identifying the target's operating system distribution and common software packages allows the attacker to select a `sudo` CVE that reliably affects the victim's environment.

---

## 2. Weaponization
The attacker builds a large compressed archive (e.g., `.tar.gz` or `.rar`) containing password-protected sub-archives alongside a trojanized Linux decryption utility (an ELF binary or wrapper script). Embedded inside the dropper is an encrypted string holding shellcode for a known `sudo` vulnerability (CVE-2025-32463 - Sudo chroot Privilege Escalation).
* **MITRE Tactic:** Resource Development ([`TA0042`](https://attack.mitre.org/tactics/TA0042/)) & Defense Evasion ([`TA0005`](https://attack.mitre.org/tactics/TA0005/))
* **MITRE Technique:** **Obtain Capabilities: Exploits** ([`T1588.005`](https://attack.mitre.org/techniques/T1588/005/)) & **Obfuscated Files or Information** ([`T1027`](https://attack.mitre.org/techniques/T1027/))
* **Explanation:** Encrypting the shellcode string inside the binary prevents static signature detection (e.g., YARA rules or antivirus scans) from flagging the exploit payload while the torrent is hosted or downloaded.

---

## 3. Delivery
The adversary seeds the malicious archive across P2P file-sharing sites and BitTorrent networks, relying on organic user downloads rather than targeted phishing.
* **MITRE Tactic:** Initial Access ([`TA0001`](https://attack.mitre.org/tactics/TA0001/))
* **MITRE Technique:** **User Execution** ([`T1204`](https://attack.mitre.org/techniques/T1204/)) (via P2P distribution)
* **Explanation:** Decentralized torrent swarms distribute the trojanized archive indiscriminately to Linux hosts without triggering perimeter email or web-gateway alerts.

---

## 4. Exploitation (User Execution & Privilege Escalation)
When the user extracts the archive and runs the bundled decryptor utility (granting it execute permissions via `chmod +x` to unlock the files), the dropper runs in user space, decrypts the embedded shellcode string in memory, and triggers the `sudo` CVE to escalate privileges from a standard user to `root`.
* **MITRE Tactic:** Execution ([`TA0002`](https://attack.mitre.org/tactics/TA0002/)) & Privilege Escalation ([`TA0004`](https://attack.mitre.org/tactics/TA0004/))
* **MITRE Techniques:**
  * **User Execution: Malicious File** ([`T1204.002`](https://attack.mitre.org/techniques/T1204/002/)): The victim executes the trojanized decryptor binary.
  * **Deobfuscate/Decode Files or Information** ([`T1140`](https://attack.mitre.org/techniques/T1140/)): The dropper decrypts the obfuscated shellcode string at runtime.
  * **Exploitation for Privilege Escalation** ([`T1068`](https://attack.mitre.org/techniques/T1068/)): The shellcode exploits the `sudo` vulnerability to bypass authentication/authorization and gain `root` execution.
* **Explanation:** Gaining `root` access is the critical pivot point in this variant; without it, the malware would remain confined to the downloading user's home directory and lack system-wide visibility.

---

## 5. Installation
Operating with `root` privileges, the dropper writes the main backdoor payload to a system directory (e.g., `/usr/local/bin/`) and registers a malicious program entry under `supervisord` so the backdoor runs automatically on boot and restarts if terminated, while the user's legitimate archive finishes decrypting.
* **MITRE Tactic:** Persistence ([`TA0003`](https://attack.mitre.org/tactics/TA0003/)) & Defense Evasion ([`TA0005`](https://attack.mitre.org/tactics/TA0005/))
* **MITRE Techniques:**
  * **Create or Modify System Process** ([`T1543`](https://attack.mitre.org/techniques/T1543/)): Dropping a `.conf` file into `/etc/supervisor/conf.d/` with `user=root`, `autostart=true`, and `autorestart=true`, then running `supervisorctl reread` and `supervisorctl update`, ensures persistent, root-level execution across reboots.
  * **Masquerading: Masquerade as Legitimate Application** ([`T1036.005`](https://attack.mitre.org/techniques/T1036/005/)): Naming the configuration file and `[program:...]` block to mimic a benign system worker (e.g., `sys-metrics-worker.conf`) hides it from casual `supervisorctl status` or `ps` inspection.

---

## 6. Command and Control (C2)
The newly spawned `supervisord`-managed daemon initiates periodic outbound beacons to the attacker's C2 infrastructure over standard network protocols.
* **MITRE Tactic:** Command and Control ([`TA0011`](https://attack.mitre.org/tactics/TA0011/))
* **MITRE Techniques:**
  * **Application Layer Protocol: Web Protocols** ([`T1071.001`](https://attack.mitre.org/techniques/T1071/001/)): Using HTTPS/HTTP or DNS to blend beacon traffic with normal Linux package updates or web requests.
  * **Scheduled Task/Job** ([`T1053`](https://attack.mitre.org/techniques/T1053/)) / Automated Beaconing: Configured timer intervals or sleep loops within the managed process control regular check-ins.

---

## 7. Actions on Objectives
Because P2P seeding infects a broad pool of machines indiscriminately, the daemon continuously harvests and transmits host telemetry to help operators identify high-value targets (such as corporate developer workstations, cloud-connected servers, or specific geographic regions).
* **MITRE Tactic:** Discovery ([`TA0007`](https://attack.mitre.org/tactics/TA0007/)) & Exfiltration ([`TA0010`](https://attack.mitre.org/tactics/TA0010/))
* **MITRE Techniques:**
  * **System Information Discovery** ([`T1082`](https://attack.mitre.org/techniques/T1082/)): Collecting OS kernel version (`uname -a`), hardware specs, CPU architecture, and installed packages.
  * **System Network Configuration Discovery** ([`T1016`](https://attack.mitre.org/techniques/T1016/)): Gathering IP addresses, routing tables, and DNS settings to determine if the host sits inside an enterprise network.
  * **Automated Exfiltration** ([`T1020`](https://attack.mitre.org/techniques/T1020/)): Packaging and transmitting the collected system profile to the C2 server on a recurring schedule.

---

## Summary Mapping Table

| Cyber Kill Chain Phase | Linux Attack Activity | MITRE ATT&CK Tactic | MITRE ATT&CK Technique |
| :--- | :--- | :--- | :--- |
| **1. Reconnaissance** | Profiling Linux P2P users and targeting vulnerable `sudo` versions | Reconnaissance ([`TA0043`](https://attack.mitre.org/tactics/TA0043/)) | Gather Victim Identity Information ([`T1589`](https://attack.mitre.org/techniques/T1589/)) |
| **2. Weaponization** | Encrypting `sudo` exploit shellcode inside a trojanized Linux decryptor | Resource Development ([`TA0042`](https://attack.mitre.org/tactics/TA0042/)) / Defense Evasion ([`TA0005`](https://attack.mitre.org/tactics/TA0005/)) | Obtain Capabilities: Exploits ([`T1588.005`](https://attack.mitre.org/techniques/T1588/005/)), Obfuscated Files ([`T1027`](https://attack.mitre.org/techniques/T1027/)) |
| **3. Delivery** | Seeding the compressed archive over BitTorrent / P2P trackers | Initial Access ([`TA0001`](https://attack.mitre.org/tactics/TA0001/)) | User Execution ([`T1204`](https://attack.mitre.org/techniques/T1204/)) |
| **4. Exploitation** | Victim runs decryptor; dropper decrypts shellcode and exploits `sudo` CVE for `root` | Execution ([`TA0002`](https://attack.mitre.org/tactics/TA0002/)) / Privilege Escalation ([`TA0004`](https://attack.mitre.org/tactics/TA0004/)) | User Execution ([`T1204.002`](https://attack.mitre.org/techniques/T1204/002/)), Deobfuscate ([`T1140`](https://attack.mitre.org/techniques/T1140/)), Exploitation for Privilege Escalation ([`T1068`](https://attack.mitre.org/techniques/T1068/)) |
| **5. Installation** | Dropper uses `root` access to register a persistent `supervisord` program | Persistence ([`TA0003`](https://attack.mitre.org/tactics/TA0003/)) / Defense Evasion ([`TA0005`](https://attack.mitre.org/tactics/TA0005/)) | Create or Modify System Process ([`T1543`](https://attack.mitre.org/techniques/T1543/)), Masquerading ([`T1036.005`](https://attack.mitre.org/techniques/T1036/005/)) |
| **6. Command & Control** | Malicious daemon establishes regular outbound C2 communications | Command and Control ([`TA0011`](https://attack.mitre.org/tactics/TA0011/)) | Application Layer Protocol: Web Protocols ([`T1071.001`](https://attack.mitre.org/techniques/T1071/001/)) |
| **7. Actions on Objectives** | Daemon regularly collects and transmits host/network specs to C2 | Discovery ([`TA0007`](https://attack.mitre.org/tactics/TA0007/)) / Exfiltration ([`TA0010`](https://attack.mitre.org/tactics/TA0010/)) | System Information Discovery ([`T1082`](https://attack.mitre.org/techniques/T1082/)), Automated Exfiltration ([`T1020`](https://attack.mitre.org/techniques/T1020/)) |
