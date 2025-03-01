sudo cp ./nexus-server.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable nexus-server.service
sudo systemctl start nexus-server.service