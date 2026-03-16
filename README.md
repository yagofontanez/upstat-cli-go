# ▲ UpStat CLI

Monitoramento de serviços UpStat direto do terminal, em tempo real.
Compatível com Linux, macOS e Windows.

---

## 🔹 Funcionalidades

- Monitoramento de status de serviços UpStat (`up` / `down`)
- Latência e uptime de cada monitor
- Dashboard em tempo real com atualização automática
- Suporte a múltiplos idiomas: Português e English
- Comandos disponíveis:
  - `start` → inicia o dashboard em tempo real
  - `logout` → remove a API key salva

---

## 🔹 Instalação

### **Linux / macOS**

O CLI é distribuído como binário via GitHub Release.
Para instalar facilmente, rode o script de instalação:

```bash
curl -sSL https://raw.githubusercontent.com/yagofontanez/upstat-cli-go/main/install.sh | bash
```

> O script detecta seu sistema, baixa o binário correto e adiciona ao PATH.

Após a instalação, use:

```bash
upstat start
```

### **Windows**

1. Baixe o binário `upstat.exe` da [última Release](https://github.com/yagofontanez/upstat-cli-go/releases)
2. Coloque o executável em uma pasta, por exemplo: `C:\upstat`
3. Adicione esta pasta ao PATH do Windows
4. Abra o PowerShell ou CMD e rode:

```powershell
upstat.exe start
```

---

## 🔹 Uso

### Iniciar dashboard

```bash
upstat start
```

Se for a primeira vez, será solicitado:

- Idioma: Português ou English
- API key do UpStat (Settings → API Keys)

### Remover API key salva

```bash
upstat logout
```

---

## 🔹 Binários

- Linux 64-bit: `upstat-linux`
- macOS 64-bit: `upstat-mac`
- Windows 64-bit: `upstat.exe`

Todos os binários estão disponíveis na seção [Releases](https://github.com/yagofontanez/upstat-cli-go/releases).

---

## 🔹 Requisitos

- Nenhum requisito extra: o CLI é **standalone**
- Para Linux/macOS, o script de instalação exige `curl`

---

## 🔹 Contribuição

PRs e issues são bem-vindos!
Para desenvolvimento local:

```bash
git clone https://github.com/yagofontanez/upstat-cli-go.git
cd upstat-cli
go run main.go
```

---

## 🔹 Licença

MIT License
