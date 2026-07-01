ARG VERSION=2.71.0
FROM mcr.microsoft.com/azure-cli:${VERSION}

ADD scripts/dump_command_table.py /

CMD ["python", "/dump_command_table.py"]
