# Windows Cyber Kill Chain & MITRE ATT&CK Breakdown (Darkhotel APT)

This attack by the **Darkhotel APT** follows a structured progression across the Lockheed Martin Cyber Kill Chain, using peer-to-peer (P2P) networks and trojanized decryption utilities to distribute the **Karba** backdoor (`Catch.exe`).

Below is the breakdown of each Kill Chain phase mapped to the corresponding **MITRE ATT&CK** Tactics, Techniques, and Procedures (TTPs) observed in the scenario.

---

## 1. Reconnaissance
Darkhotel identified the browsing habits and demographic interests of potential victims, specifically Chinese-speaking users downloading Japanese anime and manga content via P2P sharing sites.
* **MITRE Tactic:** Reconnaissance ([`TA0043`](https://attack.mitre.org/tactics/TA0043/))
* **MITRE Technique:** **Gather Victim Identity Information** ([`T1589`](https://attack.mitre.org/techniques/T1589/)) / **Search Open Websites/Domains** ([`T1593`](https://attack.mitre.org/techniques/T1593/))
* **Explanation:** Adversaries study target communities and popular P2P trends to select high-demand lure themes (such as translated manga volumes) that maximize download volume without requiring spear-phishing.

---

## 2. Weaponization
The attackers bundled a legitimate decryption utility (`AxDecrypt.exe`) together with a malicious dropper for the Karba Trojan (`Catch.exe`) inside a 900 MB RAR archive containing encrypted ZIP files and a password file.
* **MITRE Tactic:** Resource Development ([`TA0042`](https://attack.mitre.org/tactics/TA0042/))
* **MITRE Technique:** **Obtain Capabilities: Malware** ([`T1588.001`](https://attack.mitre.org/techniques/T1588/001/)) & **Stage Capabilities: Upload Malware** ([`T1608.001`](https://attack.mitre.org/techniques/T1608/001/))
* **Explanation:** Rather than relying on a software exploit, Darkhotel weaponized a functional utility by binding it with a backdoor dropper, ensuring the user receives the expected content while silently deploying malware.

---

## 3. Delivery
Darkhotel seeded the trojanized RAR archive (`(一般コミック) [古味直志] ニセコイ 第01-09巻.rar`) across Japanese P2P file-sharing sites and BitTorrent networks, achieving over 30,000 downloads.
* **MITRE Tactic:** Initial Access ([`TA0001`](https://attack.mitre.org/tactics/TA0001/))
* **MITRE Technique:** **Drive-by Compromise** ([`T1189`](https://attack.mitre.org/techniques/T1189/)) / **User Execution** ([`T1204`](https://attack.mitre.org/techniques/T1204/)) (via P2P networks)
* **Explanation:** By leveraging decentralized BitTorrent swarms and P2P portals, the adversary distributed the payload indiscriminately, allowing victims to pull the malicious archive directly onto their systems.

---

## 4. Exploitation
Because the delivery relied on a trojanized utility rather than a memory-corruption vulnerability, "exploitation" in this kill chain occurs through user-assisted execution.
* **MITRE Tactic:** Execution ([`TA0002`](https://attack.mitre.org/tactics/TA0002/))
* **MITRE Technique:** **User Execution: Malicious File** ([`T1204.002`](https://attack.mitre.org/techniques/T1204/002/))
* **Explanation:** The victim extracts the RAR archive, sees encrypted ZIP files, and runs the bundled `AxDecrypt.exe` using the provided password file to access the comics, triggering the hidden dropper.

---

## 5. Installation
When the user runs `AxDecrypt.exe`, the bound dropper extracts and installs the Karba Trojan (`Catch.exe`, detected as `Backdoor.Win32.Agent.dgrn`) onto the host system while simultaneously decrypting the legitimate files so the user suspects nothing.
* **MITRE Tactic:** Defense Evasion ([`TA0005`](https://attack.mitre.org/tactics/TA0005/)) & Persistence ([`TA0003`](https://attack.mitre.org/tactics/TA0003/))
* **MITRE Technique:** **Masquerading: Masquerade as Legitimate Application** ([`T1036.005`](https://attack.mitre.org/techniques/T1036/005/)) & **Boot or Logon Autostart Execution** ([`T1547`](https://attack.mitre.org/techniques/T1547/))
* **Explanation:** Binding the dropper to a working decryptor masks the malicious activity behind expected application behavior, while the Karba trojan establishes persistence on the Windows machine.

---

## 6. Command and Control (C2)
Once active on the victim system, `Catch.exe` initiates outbound network connections to Darkhotel's external C2 servers to register the infected host and await instructions.
* **MITRE Tactic:** Command and Control ([`TA0011`](https://attack.mitre.org/tactics/TA0011/))
* **MITRE Technique:** **Application Layer Protocol: Web Protocols** ([`T1071.001`](https://attack.mitre.org/techniques/T1071/001/))
* **Explanation:** The Karba backdoor communicates over standard web protocols (such as HTTP/HTTPS) to blend in with normal traffic, enabling operators to profile the victim machine and determine whether it belongs to a high-value target.

---

## 7. Actions on Objectives
Darkhotel historically uses Karba to triage indiscriminate P2P infections and selectively deploy second-stage espionage payloads.
* **MITRE Tactic:** Discovery ([`TA0007`](https://attack.mitre.org/tactics/TA0007/)) & Collection ([`TA0009`](https://attack.mitre.org/tactics/TA0009/))
* **MITRE Technique:** **System Information Discovery** ([`T1082`](https://attack.mitre.org/techniques/T1082/)) & **Data Staged** ([`T1074`](https://attack.mitre.org/techniques/T1074/))
* **Explanation:** Because P2P distribution infects thousands of non-target users alongside high-value executives or engineers, Darkhotel uses Karba to fingerprint infected hosts (checking system specs, installed software, and network environment) before exfiltrating sensitive data or deploying targeted modules.

---

## Summary Mapping Table

| Cyber Kill Chain Phase | Attack Activity in Scenario | MITRE ATT&CK Tactic | MITRE ATT&CK Technique |
| :--- | :--- | :--- | :--- |
| **1. Reconnaissance** | Identifying interest in Chinese-translated Japanese manga on P2P sites | Reconnaissance ([`TA0043`](https://attack.mitre.org/tactics/TA0043/)) | Gather Victim Identity Information ([`T1589`](https://attack.mitre.org/techniques/T1589/)) |
| **2. Weaponization** | Binding Karba dropper (`Catch.exe`) to `AxDecrypt.exe` inside a 900 MB RAR | Resource Development ([`TA0042`](https://attack.mitre.org/tactics/TA0042/)) | Stage Capabilities ([`T1608`](https://attack.mitre.org/techniques/T1608/)) |
| **3. Delivery** | Seeding the torrent across Japanese P2P and BitTorrent networks | Initial Access ([`TA0001`](https://attack.mitre.org/tactics/TA0001/)) | User Execution / P2P Distribution ([`T1204`](https://attack.mitre.org/techniques/T1204/)) |
| **4. Exploitation** | Victim extracts the RAR and launches `AxDecrypt.exe` to unlock the ZIPs | Execution ([`TA0002`](https://attack.mitre.org/tactics/TA0002/)) | User Execution: Malicious File ([`T1204.002`](https://attack.mitre.org/techniques/T1204/002/)) |
| **5. Installation** | Dropper silently installs `Catch.exe` while decrypting the comic files | Defense Evasion ([`TA0005`](https://attack.mitre.org/tactics/TA0005/)) / Persistence ([`TA0003`](https://attack.mitre.org/tactics/TA0003/)) | Masquerading ([`T1036.005`](https://attack.mitre.org/techniques/T1036/005/)), Boot or Logon Autostart Execution ([`T1547`](https://attack.mitre.org/techniques/T1547/)) |
| **6. Command & Control** | `Backdoor.Win32.Agent.dgrn` connects to Darkhotel C2 infrastructure | Command and Control ([`TA0011`](https://attack.mitre.org/tactics/TA0011/)) | Application Layer Protocol ([`T1071.001`](https://attack.mitre.org/techniques/T1071/001/)) |
| **7. Actions on Objectives** | Host profiling to filter high-value targets from indiscriminate victims | Discovery ([`TA0007`](https://attack.mitre.org/tactics/TA0007/)) | System Information Discovery ([`T1082`](https://attack.mitre.org/techniques/T1082/)) |
