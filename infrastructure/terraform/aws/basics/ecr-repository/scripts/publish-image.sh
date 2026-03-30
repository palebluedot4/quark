#!/usr/bin/env bash
set -euo pipefail

# aws ecr get-login-password --region <aws_region> | docker login --username AWS --password-stdin <aws_account_id>.dkr.ecr.<aws_region>.amazonaws.com
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 000000000000.dkr.ecr.ap-northeast-1.amazonaws.com

# docker buildx build \
#     --platform <platform1,platform2,...> \
#     --tag <aws_account_id>.dkr.ecr.<aws_region>.amazonaws.com/<ecr_repo>:<tag> \
#     --tag <aws_account_id>.dkr.ecr.<aws_region>.amazonaws.com/<ecr_repo>:<commit_hash> \
#     --cache-from type=registry,ref=<aws_account_id>.dkr.ecr.<aws_region>.amazonaws.com/<ecr_repo>:<cache_tag> \
#     --cache-to type=registry,ref=<aws_account_id>.dkr.ecr.<aws_region>.amazonaws.com/<ecr_repo>:<cache_tag>,mode=max \
#     --push \
#     <build_context_path>
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --tag "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/palebluedot4-dev:latest" \
    --tag "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/palebluedot4-dev:sha-$(git rev-parse --short HEAD)" \
    --cache-from "type=registry,ref=000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/palebluedot4-dev-cache:build-cache" \
    --cache-to "type=registry,ref=000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/palebluedot4-dev-cache:build-cache,mode=max" \
    --push \
    .
