FROM nvidia/cuda:13.3.1-devel-ubuntu26.04 AS builder

WORKDIR /build
RUN apt update -y && apt install git -y
RUN git clone https://github.com/ggml-org/whisper.cpp.git
RUN cmake -B build -DGGML_CUDA=1 && cmake --build build -j8 --config Release

FROM nvidia/cuda:13.3.1-runtime-ubuntu26.04
COPY --from=builder /build/bin/whisper-server /app/whisper-server
WORKDIR /app