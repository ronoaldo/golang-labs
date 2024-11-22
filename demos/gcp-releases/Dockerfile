FROM gcr.io/dataflow-templates-base/go-template-launcher-base

ARG WORKDIR=/dataflow/template
RUN mkdir -p ${WORKDIR}

COPY pipeline ${WORKDIR}/pipeline

ENV FLEX_TEMPLATE_GO_BINARY="${WORKDIR}/pipeline"

ENTRYPOINT ["/opt/google/dataflow/go_template_launcher"]
