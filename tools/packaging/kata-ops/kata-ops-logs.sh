#!/bin/bash

: ${JOURNAL_DIR:=/var/log/journal/}

NS="default"
usage() {
  cat <<EOF
Usage: $0 <pod-name> [-n namespace]

Prints the VM kernel logs corresponding to the Pod
EOF
}

get_sandbox_id() {
  journalctl -D "${JOURNAL_DIR}" -u containerd -o cat --output-fields=MESSAGE | grep "podsandboxid=" | grep "$1_$2_" | tail -1 | sed "s/^.*podsandboxid=//g"

}

main() {
  if [[ $# -ne 1 ]] && [[ $# -ne 3 ]]; then
    usage
    exit 1
  fi
  if [[ $# -eq 3 ]] && [[ "$2" != "-n"  ]]; then
    usage
    exit 1
  fi

  if [[ $# -eq 3 ]]; then
    export NS="$3"
  fi

  echo "Getting sanbox id for $1 in namespace $NS..."
  SBID=$(get_sandbox_id $1 $NS)
  if [[ "$SBID" == "" ]]; then
    echo "Sandbox for Pod $1 in namespace $NS not found"
    exit 1
  fi

  echo "Getting logs for sandbox id $SBID..."
  echo ""
  journalctl -D "${JOURNAL_DIR}" -u containerd -o cat | grep "sandbox=$SBID" | grep "vmconsole=\"\[" | sed "s/^.*vmconsole=//g"
}

main "$@"
