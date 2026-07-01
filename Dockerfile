ARG VERSION=2.71.0
FROM mcr.microsoft.com/azure-cli:${VERSION}

COPY scripts/dump_command_table.py /

CMD ["python3", "/dump_command_table.py"]
