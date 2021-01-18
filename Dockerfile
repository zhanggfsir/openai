FROM golang

COPY openai-backend .
COPY conf ./conf

CMD [ "./openai-backend" ]