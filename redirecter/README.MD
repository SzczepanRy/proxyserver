windows divert-go
linux go-nfqueue



#linux required installs :
  sudo pacman -S libnetfilter_queue

  #czyszczenie ip tables
  asudo iptables -F

  sudo setcap 'cap_net_admin,cap_net_raw+ep' redirecter_linux
  sudo iptables -I OUTPUT -p tcp -j NFQUEUE --queue-num 100


  # 1. Ignoruj ruch lokalny (localhost) - to zapobiegnie zawieszeniu usług systemowych
    sudo iptables -A OUTPUT -o lo -j ACCEPT
    sudo iptables -A INPUT -i lo -j ACCEPT

  # 2. WYJĄTEK DLA SSH (Port 22) - KONIECZNE jeśli łączysz się zdalnie (np. Putty/Termius)
    sudo iptables -A OUTPUT -p tcp --sport 22 -j ACCEPT
    sudo iptables -A OUTPUT -p tcp --dport 22 -j ACCEPT
    sudo iptables -A INPUT -p tcp --sport 22 -j ACCEPT
    sudo iptables -A INPUT -p tcp --dport 22 -j ACCEPT

  # 3. WYJĄTEK DLA TWOJEGO PROGRAMU - Zapobiega pętli (zakładam, że odpalasz przez sudo, więc user to root)
    sudo iptables -A OUTPUT -p tcp -m owner --uid-owner root -j ACCEPT

  # 4. GŁÓWNA REGUŁA - Skieruj CAŁY pozostały ruch TCP do Twojego programu (kolejka 10)
    sudo iptables -A OUTPUT -p tcp -j NFQUEUE --queue-num 10
    sudo iptables -A INPUT -p tcp -j NFQUEUE --queue-num 10

