ARG VERSION=2.71.0
FROM mcr.microsoft.com/azure-cli:${VERSION}

COPY scripts/dump_command_table.py /

ENV PYTHONPATH=/usr/lib64/az/lib/python3.12/site-packages:/usr/lib/az/lib/python3.12/site-packages

ENTRYPOINT []
CMD ["python3", "/dump_command_table.py"]
