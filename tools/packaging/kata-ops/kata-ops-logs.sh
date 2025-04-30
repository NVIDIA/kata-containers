#!/bin/bash

: ${JOURNAL_DIR:=/var/log/journal/}

NS="default"
FOLLOW="false"
usage() {
  cat <<EOF
Usage: $0 <pod-name> [-n namespace] [--follow]

Prints the VM kernel logs corresponding to the Pod
EOF
}

get_sandbox_id() {
  journalctl -D "${JOURNAL_DIR}" -u containerd -o cat --output-fields=MESSAGE | grep "podsandboxid=" | grep "$1_$2_" | tail -1 | sed "s/^.*podsandboxid=//g"

}

main() {
  if [[ $# -lt 1 ]] || [[ $# -gt 4 ]]; then
    usage
    exit 1
  fi
  if [[ $# -eq 3 ]] && [[ "$2" != "-n"  ]]; then
    usage
    exit 1
  fi

  if [[ $# -eq 3 ]] && [[ "$2" == "-n" ]]; then
    export NS="$3"
  fi

  if [[ $# -eq 2 ]] && [[ "$2" == "--follow" ]]; then
    export FOLLOW="true"
  fi

  if [[ $# -eq 4 ]] && [[ "$4" == "--follow" ]]; then
    export FOLLOW="true"
  fi

  echo "Getting sanbox id for $1 in namespace $NS..."
  SBID=$(get_sandbox_id $1 $NS)
  if [[ "$SBID" == "" ]]; then
    echo "Sandbox for Pod $1 in namespace $NS not found"
    exit 1
  fi

  echo "Getting logs for sandbox id $SBID..."
  echo ""
  if [[ $FOLLOW == "true" ]]; then
    journalctl -D "${JOURNAL_DIR}" -f --no-tail -u containerd -o cat | grep "sandbox=$SBID" | grep "vmconsole=\"\[" | sed "s/^.*vmconsole=//g"
  else
    journalctl -D "${JOURNAL_DIR}" -u containerd -o cat | grep "sandbox=$SBID" | grep "vmconsole=\"\[" | sed "s/^.*vmconsole=//g"
  fi
}

main "$@"
