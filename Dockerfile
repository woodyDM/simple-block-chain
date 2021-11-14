FROM centos:8
WORKDIR /app
COPY . /app/
COPY application-prod.json /app/
RUN rm -f /etc/localtime \
&& ln -sv /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone

EXPOSE 80/tcp
CMD /app/app -p prod